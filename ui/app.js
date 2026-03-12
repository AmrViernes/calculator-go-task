// Points to our backend API (same host since we serve both together)
const API_BASE = window.location.origin;

// Grab all the DOM elements we'll need upfront
const calculateForm = document.getElementById('calculateForm');
const configForm = document.getElementById('configForm');
const orderQuantityInput = document.getElementById('orderQuantity');
const resultSection = document.getElementById('resultSection');
const errorSection = document.getElementById('errorSection');
const errorMessage = document.getElementById('errorMessage');
const orderDisplay = document.getElementById('orderDisplay');
const totalItemsDisplay = document.getElementById('totalItemsDisplay');
const totalPacksDisplay = document.getElementById('totalPacksDisplay');
const packsList = document.getElementById('packsList');
const packSizesContainer = document.getElementById('packSizesContainer');
const addPackSizeBtn = document.getElementById('addPackSize');

// Boot up the app when the page loads
document.addEventListener('DOMContentLoaded', () => {
	loadPackSizes();
	setupEventListeners();
});

// Wire up all our event handlers
function setupEventListeners() {
	// When someone hits calculate on the form
	calculateForm.addEventListener('submit', handleCalculate);

	// When someone saves new pack sizes
	configForm.addEventListener('submit', handleUpdateConfig);

	// Add a new pack size input field
	addPackSizeBtn.addEventListener('click', addPackSizeInput);

	// Quick test buttons - just click to auto-fill and calculate
	document.querySelectorAll('.test-buttons .btn').forEach(btn => {
		btn.addEventListener('click', () => {
			const quantity = parseInt(btn.dataset.quantity);
			orderQuantityInput.value = quantity;
			calculatePacks(quantity);
		});
	});

	// Preset buttons for loading pre-defined pack configs
	document.querySelectorAll('.preset-buttons .btn').forEach(btn => {
		btn.addEventListener('click', () => {
			const sizes = JSON.parse(btn.dataset.sizes);
			updatePackSizes(sizes);
		});
	});
}

// Form was submitted - validate and call the API
async function handleCalculate(e) {
	e.preventDefault();
	const quantity = parseInt(orderQuantityInput.value);
	if (isNaN(quantity) || quantity < 0) {
		showError('Please enter a valid non-negative number');
		return;
	}
	await calculatePacks(quantity);
}

// Call the backend to calculate optimal packs
async function calculatePacks(quantity) {
	try {
		hideError();
		hideResult();

		const response = await fetch(`${API_BASE}/api/calculate`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ orderQuantity: quantity }),
		});

		if (!response.ok) {
			const error = await response.text();
			throw new Error(error || 'Failed to calculate packs');
		}

		const data = await response.json();
		displayResult(data, quantity);
	} catch (error) {
		showError(error.message || 'An error occurred while calculating');
	}
}

// Show the calculated packs to the user
function displayResult(data, orderQuantity) {
	orderDisplay.textContent = orderQuantity.toLocaleString();

	let totalItems = 0;
	let totalPacks = 0;

	// Build the pack display elements
	packsList.innerHTML = '';
	data.packs.forEach(pack => {
		totalItems += pack.size * pack.count;
		totalPacks += pack.count;

		const li = document.createElement('li');
		li.className = 'pack-item';
		li.innerHTML = `
			<span class="count">${pack.count}</span>
			<span class="size">× ${pack.size.toLocaleString()}</span>
			<span class="label">pack${pack.count > 1 ? 's' : ''}</span>
		`;
		packsList.appendChild(li);
	});

	totalItemsDisplay.textContent = totalItems.toLocaleString();
	totalPacksDisplay.textContent = totalPacks;

	showResult();
}

// Get the current pack sizes from the server
async function loadPackSizes() {
	try {
		const response = await fetch(`${API_BASE}/api/packsizes`);
		if (!response.ok) {
			throw new Error('Failed to load pack sizes');
		}
		const data = await response.json();
		renderPackSizeInputs(data.packSizes);
	} catch (error) {
		console.error('Failed to load pack sizes:', error);
		// Fall back to defaults if the API is down
		renderPackSizeInputs([250, 500, 1000, 2000, 5000]);
	}
}

// Draw the pack size input fields
function renderPackSizeInputs(sizes) {
	packSizesContainer.innerHTML = '';
	sizes.forEach(size => addPackSizeInput(size));
}

// Add a new pack size input row
function addPackSizeInput(value = '') {
	const div = document.createElement('div');
	div.className = 'pack-size-input';
	div.innerHTML = `
		<input type="number" min="1" value="${value}" placeholder="Size" required>
		<span>items</span>
		<button type="button" class="btn btn-danger remove-pack" aria-label="Remove pack size">×</button>
	`;

	// Wire up the remove button
	div.querySelector('.remove-pack').addEventListener('click', () => {
		if (packSizesContainer.children.length > 1) {
			div.remove();
		} else {
			showError('At least one pack size is required');
		}
	});

	packSizesContainer.appendChild(div);
}

// Save the new pack sizes to the server
async function handleUpdateConfig(e) {
	e.preventDefault();

	// Gather all the input values
	const sizes = [];
	const inputs = packSizesContainer.querySelectorAll('input');

	for (const input of inputs) {
		const value = parseInt(input.value);
		if (isNaN(value) || value <= 0) {
			showError('All pack sizes must be positive numbers');
			return;
		}
		sizes.push(value);
	}

	if (sizes.length === 0) {
		showError('At least one pack size is required');
		return;
	}

	await updatePackSizes(sizes);
}

// Send the new pack sizes to the API
async function updatePackSizes(sizes) {
	try {
		hideError();

		const response = await fetch(`${API_BASE}/api/packsizes`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ packSizes: sizes }),
		});

		if (!response.ok) {
			const error = await response.text();
			throw new Error(error || 'Failed to update pack sizes');
		}

		const data = await response.json();
		renderPackSizeInputs(data.packSizes);
		hideResult();

		// Re-calculate with the new pack sizes if there's a value in the input
		const currentQuantity = parseInt(orderQuantityInput.value);
		if (currentQuantity > 0) {
			await calculatePacks(currentQuantity);
		}
	} catch (error) {
		showError(error.message || 'An error occurred while updating pack sizes');
	}
}

// Show/hide result section
function showResult() {
	resultSection.classList.remove('hidden');
}

function hideResult() {
	resultSection.classList.add('hidden');
}

// Show/hide error messages
function showError(message) {
	errorMessage.textContent = message;
	errorSection.classList.remove('hidden');
}

function hideError() {
	errorSection.classList.add('hidden');
}

// Don't let people enter negative numbers
orderQuantityInput.addEventListener('input', function() {
	if (this.value < 0) {
		this.value = 0;
	}
});

// Validate pack size inputs as they're typed
packSizesContainer.addEventListener('input', function(e) {
	if (e.target.type === 'number') {
		if (e.target.value < 1) {
			e.target.value = 1;
		}
	}
});
