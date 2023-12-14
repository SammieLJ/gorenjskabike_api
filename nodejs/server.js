const express = require('express');
const aggregateData = require('./proxy');
const app = express();
const port = 3001;

app.use(express.json());

app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS, PUT, PATCH, DELETE');
  res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization');
  next();
});

app.get('/api/localdata', async (req, res) => {
  try {
    const { selectedDateTimeFrom, selectedDateTimeTo } = req.query;

    // Validate the presence of both start and end dates
    if (!selectedDateTimeFrom || !selectedDateTimeTo) {
      return res.status(400).json({ error: 'Missing date range parameters' });
    }

    const allData = await aggregateData(selectedDateTimeFrom, selectedDateTimeTo);

    res.json(allData);
  } catch (error) {
    console.error('Error fetching and aggregating data:', error);
    res.status(500).json({ error: 'Internal Server Error' });
  }
});

app.listen(port, () => {
  console.log(`Server is running at http://localhost:${port}`);
});
