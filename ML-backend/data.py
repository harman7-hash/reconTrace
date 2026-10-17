import psycopg2
import pandas as pd
import os 

def get_data_train(user_id):
    conn = psycopg2.connect(os.getenv("CONNECTION"))
    cursor = conn.cursor()
# use the cursor to interact with your database
    cursor.execute(
    "SELECT * FROM server_metrics WHERE user_id = %s ORDER BY time ASC",
    (user_id,))
    col_names = [desc[0] for desc in cursor.description]
    data=cursor.fetchall()
    df=pd.DataFrame(data,columns=col_names)
    print(f"Data shape: {df.shape}")
    return df