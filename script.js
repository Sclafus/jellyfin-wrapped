
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
		displayJson(activity);
	} catch (error) {
		console.error('Error fetching activity:', error);
	}
}

function displayJson(data) {
	const output = document.getElementById('output');
	output.textContent = JSON.stringify(data, null, 4);
}
