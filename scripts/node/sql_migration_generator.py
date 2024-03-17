import os
import datetime
from os.path import join, dirname
from dotenv import load_dotenv
import psycopg2

dotenv_path = join(dirname(__file__), "../../.env")
load_dotenv(dotenv_path)

username = os.environ.get("SQL_USERNAME")
password = os.environ.get("SQL_PASSWORD")

if username==None or password==None:
    print("Couldn't find credentials in .env file.")
    exit(1)

conn = psycopg2.connect(
    database="nexus",
    host="localhost", 
    user=username, 
    password=password, 
    port="5432"
)

cursor = conn.cursor()

images_path = join(dirname(__file__), "../../static/images/")

trial_query = """INSERT INTO trials(trial_num) VALUES (%s)"""
image_query = """INSERT INTO images(file_num, trial_num, time) VALUES(%s,%s,%s)"""
# TODO!: ERROR HANDLING HERE TO MAKE SURE SHIT DOESN'T GO WRONG!
for trial_directory in os.listdir(images_path):
    trial_num = int(trial_directory[6:]) # first six characters are "trial-"

    cursor.execute(trial_query, (trial_num,))
    for image_path in os.listdir(join(images_path, trial_directory)):
        image_num = int(image_path[:-4]) # last four characters are ".jpg"


        path = join(images_path, trial_directory, image_path)
        m_time = os.path.getmtime(path)
        date_modified = datetime.datetime.fromtimestamp(m_time)
        db_formatted_timestamp = date_modified.strftime("%Y-%m-%d %H:%M:%S")

        cursor.execute(image_query, (image_num, trial_num, db_formatted_timestamp,))

        if(image_num==1):
            cursor.execute("UPDATE trials SET start_time = %s WHERE trial_num = %s", (db_formatted_timestamp, trial_num))


conn.commit()
cursor.close()
conn.close()

