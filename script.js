(() => {
	const yearSelect = document.getElementById('yearSelect');
	const currentYear = new Date().getFullYear();
	const startYear = 2023; // TODO: this is completely arbitrary ¯\_(ツ)_/¯

	// Clear existing options
	yearSelect.innerHTML = '<option value="All">All Years</option>';

	// Add years from current year backwards
	for (let year = currentYear; year >= startYear; year--) {
		const option = document.createElement('option');
		option.value = year;
		option.textContent = year;
		yearSelect.appendChild(option);
	}
})();

let selectedUserId = '';

async function fetchUsers() {
	const baseUrl = document.getElementById('baseUrl').value;
	const apiKey = document.getElementById('apiKey').value;
	const usersUrl = `${baseUrl}/Users`;

	try {
		//  const response = await fetch(usersUrl);
		const response = await fetch(usersUrl, {
			headers: {
				'X-Emby-Token': apiKey
			}
		});
		const users = await response.json();

		const select = document.getElementById('usersSelect');
		select.innerHTML = '<option value="">Select a user...</option>';

		users.forEach(user => {
			const option = document.createElement('option');
			option.value = user.Id;
			option.textContent = user.Name;
			select.appendChild(option);
		});
	} catch (error) {
		console.error('Error fetching users:', error);
	}
}

async function fetchActivity() {
	const baseUrl = document.getElementById('baseUrl').value;
	const apiKey = document.getElementById('apiKey').value;
	const userId = document.getElementById('usersSelect').value;
	const year = document.getElementById('yearSelect').value;

	if (!userId) {
		alert('Please select a user first!');
		return;
	}

	const activityUrl = `${baseUrl}/Items?userId=${userId}&isPlayed=true&recursive=true`;

	try {
		const response = await fetch(activityUrl, {
			headers: {
				'X-Emby-Token': apiKey
			}
		});

		const activity = await response.json();
		const filteredActivity = filterActivity(activity, year);
		console.log(filteredActivity);
		displayData(filteredActivity);
	} catch (error) {
		console.error('Error fetching activity:', error);
	}
}

function filterActivity(activity, year) {
	if (year === 'All') {
		return activity.Items;
	}
	try {
		year = parseInt(year);

		const filteredItems = activity.Items.filter(item => {
			const lastPlayedDate = item.UserData.LastPlayedDate;
			const lastPlayedYear = parseInt(lastPlayedDate.substring(0, 4));
			return lastPlayedYear === year;
		});
		return filteredItems;
	} catch (error) {
		console.error('Error filtering activity:', error);
	}
}

function displayData(data) {
	// Calculate statistics
	const seriesWatchTime = {};
	let totalMoviesTime = 0;
	let totalMoviesCount = 0;

	data.forEach(item => {
		if (item.Type === 'Episode') {
			// For TV series
			const seriesName = item.SeriesName;
			const runtimeSeconds = item.RunTimeTicks / 10000000;

			if (!seriesWatchTime[seriesName]) {
				seriesWatchTime[seriesName] = {
					name: seriesName,
					totalWatchTime: 0,
					count: 0
				};
			}
			seriesWatchTime[seriesName].totalWatchTime += runtimeSeconds;
			seriesWatchTime[seriesName].count++;
		} else if (item.Type === 'Movie') {
			// For movies
			totalMoviesTime += item.RunTimeTicks / 10000000;
			totalMoviesCount++;
		}
	});

	// Convert to array and sort by watch time in descending order
	const sortedSeries = Object.values(seriesWatchTime).sort((a, b) => b.totalWatchTime - a.totalWatchTime);
	const top10Series = sortedSeries.slice(0, 10);

	// Convert seconds to hours for display
	top10Series.forEach(series => {
		series.totalWatchTime = (series.totalWatchTime / 3600).toFixed(2);
	});

	// Format output
	const output = document.getElementById('output');
	let result = `
				<h3>Top 10 Most Watched TV Series</h3>
				<ul>`;

	top10Series.forEach(series => {
		result += `
						<li>
								${series.name}
								(<strong>${series.count} episodes</strong>,
								 <strong>${series.totalWatchTime} hours</strong>)
						</li>`;
	});

	result += `</ul>
								<h3>Movies</h3>
								<p>
										Total movies watched: ${totalMoviesCount}<br>
										Total movie watch time: ${(totalMoviesTime / 3600).toFixed(2)} hours
								</p>`;

	output.innerHTML = result;
}
