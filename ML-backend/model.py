import torch
import torch.nn as nn
from torch.utils.data import Dataset,DataLoader
import numpy as np
from config import BATCH_SIZE,HIDDEN_SIZE,INPUT_SIZE,SEQ_LEN,INITIAL_EPOCHS ,LR,FINE_EPOCHS
import torch.onnx
import time
import io
from data import get_data_train
class LSTMencoder(nn.Module):
    def __init__(self,input_size=INPUT_SIZE,hidden_size=32):
        super().__init__()
        self.lstm=nn.LSTM(input_size=input_size,hidden_size=hidden_size,batch_first=True)
    def forward(self,x):
        _,(hidden,cell_state)=self.lstm(x)
        return hidden,cell_state
class LSTMdecoder(nn.Module):
    def __init__(self,input_size=INPUT_SIZE,hidden_size=32):
        super().__init__()
        self.lstm=nn.LSTM(input_size=input_size,hidden_size=hidden_size,batch_first=True)
        self.fc=nn.Linear(hidden_size,input_size)
    def forward(self,x,hidden_state,cell_state):
        out,(hidden,cell_state)=self.lstm(x,(hidden_state,cell_state))
        prediction=self.fc(out)
        return prediction,hidden,cell_state
class LSTMautoencoder(nn.Module):
    def __init__(self,input_size=INPUT_SIZE,hidden_size=32,seq_len=SEQ_LEN):
        super().__init__()
        self.seq_len=seq_len
        self.encoder=LSTMencoder(input_size,hidden_size)
        self.decoder=LSTMdecoder(input_size,hidden_size)
        self.input_size=input_size
    def forward(self,training_seq):
        batch_size=training_seq.size(0)
        hidden_state,cell_state=self.encoder(training_seq)
        decoder_input=torch.zeros(batch_size,1,self.input_size)
        outputs=[]
        for i in range(self.seq_len):
            prediction,hidden_state,cell_state=self.decoder(decoder_input,hidden_state,cell_state)
            outputs.append(prediction)
            decoder_input=prediction
        final_outputs = torch.cat(outputs, dim=1)
        return final_outputs
class customDataset(Dataset):
    def __init__(self,df,seq_len):
        self.seq_len=seq_len
        super().__init__()
        features_df=df.drop(columns=['time','user_id'])#############################TEMPCODE#########################
        cols = features_df.iloc[:, 0:4] 
        features_df.iloc[:, 0:4] = (cols - cols.min()) / (cols.max() - cols.min())
        #########################################################
        print(f"Features shape: {features_df.shape}")
        self.dt=torch.tensor(features_df.iloc[:,:].values,dtype=torch.float32)
    def __len__(self):
        return max(0,len(self.dt)-self.seq_len+1)
    def __getitem__(self, i):
        return self.dt[i:i+self.seq_len]
def train_model(user_id):
    model=LSTMautoencoder()
    model.train()
    dataset=customDataset(get_data_train(user_id),SEQ_LEN)
    dataloader=DataLoader(dataset,batch_size=BATCH_SIZE,shuffle=True)
    loss_function=nn.MSELoss(reduction='none')
    optimizer=torch.optim.Adam(model.parameters(),lr=LR)
    dummy_input=torch.zeros(1,SEQ_LEN,INPUT_SIZE)
    for epoch in range(INITIAL_EPOCHS):
        total_epoch_loss=0
        total_traversed=0
        for data in dataloader:
            inputs=data
            preds=model(inputs)
            optimizer.zero_grad()
            loss=loss_function(preds,inputs)
            loss=loss.mean()
            loss.backward()
            optimizer.step()
            total_epoch_loss+=loss.item()*len(data)
            total_traversed+=len(data)
        train_losses = []
        model.eval()
        with torch.no_grad():
            for data in dataloader:
                inputs = data
                preds = model(inputs)
                loss = loss_function(preds, inputs)
                # Per-sequence reconstruction error
                sequence_loss = loss.mean(dim=(1, 2))
                # Store every sample's error
                train_losses.extend(sequence_loss.cpu().numpy())
            THRESHOLD = np.mean(train_losses) + 3 * np.std(train_losses)
            filename=f"{user_id}_{int(time.time())}.onnx"
            temp_file_path = f"./public/temp/{filename}"
            torch.onnx.export(model,dummy_input,temp_file_path,export_params=True, external_data=False)
            return THRESHOLD,model,temp_file_path,filename
def fine_tune_model(model,user_id):
    model.train()
    dataset=customDataset(get_data_train(user_id),SEQ_LEN)
    dataloader=DataLoader(dataset,batch_size=BATCH_SIZE,shuffle=True)
    loss_function=nn.MSELoss(reduction='none')
    optimizer=torch.optim.Adam(model.parameters(),lr=LR)
    for epoch in range(FINE_EPOCHS):
        total_epoch_loss=0
        total_traversed=0
        for data in dataloader:
            inputs=data
            preds=model(inputs)
            optimizer.zero_grad()
            loss=loss_function(preds,inputs)
            loss=loss.mean()
            loss.backward()
            optimizer.step()
            total_epoch_loss+=loss.item()*len(data)
            total_traversed+=len(data)
        train_losses = []
        model.eval()
        with torch.no_grad():
            for data in dataloader:
                inputs = data
                preds = model(inputs)
                loss = loss_function(preds, inputs)
                sequence_loss = loss.mean(dim=(1, 2))
                train_losses.extend(sequence_loss.cpu().numpy())
            THRESHOLD = np.mean(train_losses) + 3 * np.std(train_losses)
            return THRESHOLD,model