# server.py
from flask import Flask, request, jsonify
from proxy import aggregate_data

app = Flask(__name__)
port = 3001

app.config['CORS_HEADERS'] = 'Content-Type'

@app.route('/')
def index():
    return 'Hello, this is the server!'

@app.after_request
def after_request(response):
    response.headers['Access-Control-Allow-Origin'] = '*'
    response.headers['Access-Control-Allow-Methods'] = 'GET, POST, OPTIONS, PUT, PATCH, DELETE'
    response.headers['Access-Control-Allow-Headers'] = 'Content-Type, Authorization'
    return response

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
    app.run(port=port)
