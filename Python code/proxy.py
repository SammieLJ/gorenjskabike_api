# proxy.py
from flask import Flask, jsonify, request
import json
import requests
import os
import random
from datetime import datetime, timedelta, timezone

app = Flask(__name__)
app.debug = True
port = 3000

def generate_predictive_data(id, timestamp):
    return {
        'id': id,
        'lat': '46.355340188965926',
        'lon': '0',
        'name': setName_AccordiglyToID(id),
        'numberOfBikes': random.randint(4, 7),
        'numberOfElectricBikes': random.randint(1, 3),
        'numberOfFreeBikes': random.randint(4, 7),
        'numberOfFreeLocks': random.randint(4, 7),
        'numberOfLocks': '10',
        'numberOfTotalFaulty': random.choices([0, 1, 2], cum_weights=(60, 30, 10), k=1),
        'street': setStreet_AccordiglyToID(id),
        'timestamp': timestamp,
    }

def setName_AccordiglyToID(id):
    if id == '96':
        return 'Tržič - Deteljica'
    if id == '97':
        return 'Tržič - Občina Tržič'
    if id == '98':
        return 'Tržič - BPT'
    
def setStreet_AccordiglyToID(id):
    if id == '96':
        return 'Deteljica 2 3/4 P!'
    if id == '97':
        return 'Trg Svobode 18 3/3 P!'
    if id == '98':
        return 'Predilniška cesta 14 3/4 P!'

def aggregate_data(selected_datetime_from, selected_datetime_to):
    try:
        #data = None
        data = []  # Initialize data as an empty list

        current_predictive_data = None
        current_predictive_data_ID96 = None
        current_predictive_data_ID97 = None
        current_predictive_data_ID98 = None
        future_timestamp = None

        # Check if allData.json exists
        if os.path.exists('allData.json'):
            # If the file exists, read its content
            with open('allData.json', 'r', encoding='utf-8') as file:
                data = json.load(file)
        else:
            # If the file doesn't exist, make the API call with pagination
            api_url = 'https://trzic.musiclab.si/api/gorenjskabike'
            #data = []

            for page in range(1, 1016):  # Assuming a maximum of 1015 pages
                params = {'page': page, 'size': 1000}
                response = requests.get(api_url, params=params)

                try:
                    page_data = response.json()
                    if not page_data:
                        break  # No more data available
                    data.extend(page_data)
                except json.JSONDecodeError:
                    print(f'Error: Unable to parse API response as JSON for page {page}')

            # Save the aggregated data to the file
            with open('allData.json', 'w', encoding='utf-8') as file:
                json.dump(data, file, ensure_ascii=False)

        # Filter data based on the provided date range
        if selected_datetime_from and selected_datetime_to:
            start_date = selected_datetime_from
            end_date = selected_datetime_to

            # print("Start Date", start_date)
            # print("End Date", end_date)
            end_date_timestamp = datetime.fromisoformat(end_date)

            # Bike MAP spletna stran ima en date picker in sta oba enaka, treba dodat 1 min v naprej
            if start_date == end_date:
                end_date_timestamp = end_date_timestamp + timedelta(minutes=1)
                end_date = end_date_timestamp.isoformat()
            print("Start Date", start_date)
            print("End Date", end_date)
            print("End date timestamp and type:", end_date_timestamp, type(end_date_timestamp))

            # Filter data for the selected date range
            data = [entry for entry in data if start_date <= entry.get('timestamp', '') <= end_date]
            #print(data)
            print("Število zapisov v data objektu: ", len(data))

            # hard coded last api entry date 2023-10-01T01:59,  prej last_data_timestamp
            hard_coded_last_api_entry_date = datetime.fromisoformat("2023-10-01T01:59")
            print("Hard coded last api entry date and type:", hard_coded_last_api_entry_date, type(hard_coded_last_api_entry_date))

            if end_date_timestamp >= hard_coded_last_api_entry_date:
                # Generate predictive data for missing timestamps
                for timestamp in generate_missing_timestamps(start_date, end_date):
                    for id in ['96', '97', '98']:

                        # Check if there is any entry with the same ID and timestamp
                        matching_entry = next((entry for entry in data if entry['id'] == id and are_timestamps_equal(timestamp, entry['timestamp'])), None)

                        if matching_entry is None:
                            # if it is first time generating pred. data or second time that future timestamp is equal to current timestamp
                            if future_timestamp is None or are_timestamps_equal(timestamp, future_timestamp):
                                #This adds for ID 96
                                current_predictive_data_ID96 = generate_predictive_data(id, timestamp)
                                #current_predictive_data = current_predictive_data_ID96

                                current_predictive_data_ID97 = None
                                current_predictive_data_ID98 = None

                                # Randomly add minutes to the timestamp, people traveling by bikes 10 .. 40 min
                                additional_minutes = random.randint(10, 40)  # Adjust the range as needed
                                future_timestamp = add_minutes_to_timestamp(timestamp, additional_minutes)
                            else:
                                # Add for ID 97 and ID 98
                                if id == '97' and current_predictive_data_ID97 == None:
                                    current_predictive_data_ID97 = generate_predictive_data(id, timestamp)
                                    #current_predictive_data = current_predictive_data_ID97
                                if id == '98' and current_predictive_data_ID98 == None:
                                    current_predictive_data_ID98 = generate_predictive_data(id, timestamp)
                                    #current_predictive_data = current_predictive_data_ID98

                            # for looping in between future_timestamp and current timestamp, for ID chnage to current predictive data
                            if id == '96':
                                current_predictive_data = current_predictive_data_ID96
                            if id == '97':
                                current_predictive_data = current_predictive_data_ID97
                            if id == '98':
                                current_predictive_data = current_predictive_data_ID98

                            data.append(current_predictive_data)
        print("Število zapisov v data (na konc): ", len(data))
        return data
    except Exception as e:
        print('Error fetching and aggregating data:', e)
        return []

def are_timestamps_equal(timestamp1, timestamp2):
    # Parse the timestamps to datetime objects
    #dt1 = datetime.fromisoformat(timestamp1)
    #dt2 = datetime.fromisoformat(timestamp2)
    dt1 = datetime.fromisoformat(trimDateStr(timestamp1))
    dt2 = datetime.fromisoformat(trimDateStr(timestamp2))

    # Truncate microseconds to seconds
    dt1 = dt1.replace(microsecond=0)
    dt2 = dt2.replace(microsecond=0)

    # Convert both timestamps to UTC
    dt1_utc = dt1.replace(tzinfo=timezone.utc)
    dt2_utc = dt2.replace(tzinfo=timezone.utc)

    # Check if the UTC timestamps are equal
    return dt1_utc == dt2_utc

def trimDateStr(timestamp_str):
    const_ts_len = 19
    timestamp_str_returned = ''
    if len(timestamp_str) < const_ts_len:
         timestamp_str_returned = timestamp_str
    else:
        timestamp_str_returned = timestamp_str[:const_ts_len]

    #print('Desno trimmed TS: ', timestamp_str_returned)
    return timestamp_str_returned
def add_minutes_to_timestamp(timestamp, minutes):
    dt = datetime.fromisoformat(timestamp)
    dt_with_additional_minutes = dt + timedelta(minutes=minutes)
    return dt_with_additional_minutes.isoformat()

def generate_missing_timestamps(start_date, end_date):
    current_date = datetime.fromisoformat(start_date)
    end_date = datetime.fromisoformat(end_date)
    while current_date <= end_date:
        yield current_date.isoformat()
        current_date += timedelta(minutes=1)

@app.route('/')
def index():
    return 'Hello, this is the proxy server!'

@app.route('/api/localdata', methods=['GET'])
def get_local_data():
    try:
        selected_datetime_from = request.args.get('selectedDateTimeFrom')
        selected_datetime_to = request.args.get('selectedDateTimeTo')

        # Validate the presence of both start and end dates
        if not selected_datetime_from or not selected_datetime_to:
            return jsonify({'error': 'Missing date range parameters'}), 400

        all_data = aggregate_data(selected_datetime_from, selected_datetime_to)

        return jsonify(all_data)
    except Exception as e:
        print('Error fetching and aggregating data:', e)
        return jsonify({'error': 'Internal Server Error'}), 500

if __name__ == '__main__':
    app.run(port=port, debug=True)
