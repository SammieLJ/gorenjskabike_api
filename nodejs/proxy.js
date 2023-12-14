// proxy.js
const express = require('express');
const fs = require('fs');

const app = express();
const port = 3000;

const aggregateData = async (selectedDateTimeFrom, selectedDateTimeTo) => {
  try {
    // Use dynamic import for node-fetch
    const fetch = (await import('node-fetch')).default;

    let data;

    // Check if allData.json exists
    if (fs.existsSync('allData.json')) {
      // If the file exists, read its content
      data = JSON.parse(fs.readFileSync('allData.json'));
    } else {
      // If the file doesn't exist, make the API call
      const apiURL = 'https://trzic.musiclab.si/api/gorenjskabike';
      const response = await fetch(apiURL);
      data = await response.json();

      // Save the API response to the file
      fs.writeFileSync('allData.json', JSON.stringify(data));
    }

    // Filter data based on the provided date range
    if (selectedDateTimeFrom && selectedDateTimeTo) {
      const startDate = new Date(selectedDateTimeFrom);
      const endDate = new Date(selectedDateTimeTo);

      // Filter data for the selected date range
      data = data.filter(entry => {
        const entryDate = new Date(entry.timestamp);
        return entryDate >= startDate && entryDate <= endDate;
      });
    }

    return data;
  } catch (error) {
    console.error('Error fetching and aggregating data:', error);
    throw error;
  }
};

module.exports = aggregateData;

app.listen(port, () => {
  console.log(`Server is running on port ${port}`);
});
