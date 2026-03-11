// API Base URL - adjust if API is on different port/host
const API_BASE = window.location.origin;

// DOM Elements
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

// Initialize the app
document.addEventListener('DOMContentLoaded', () => {
    loadPackSizes();
    setupEventListeners();
});

// Setup Event Listeners
function setupEventListeners() {
    // Calculate form submission
    calculateForm.addEventListener('submit', handleCalculate);

    // Config form submission
    configForm.addEventListener('submit', handleUpdateConfig);

    // Add pack size button
    addPackSizeBtn.addEventListener('click', addPackSizeInput);

    // Test case buttons
    document.querySelectorAll('.test-buttons .btn').forEach(btn => {
        btn.addEventListener('click', () => {
            const quantity = parseInt(btn.dataset.quantity);
            orderQuantityInput.value = quantity;
            calculatePacks(quantity);
        });
    });

    // Preset buttons
    document.querySelectorAll('.preset-buttons .btn').forEach(btn => {
        btn.addEventListener('click', () => {
            const sizes = JSON.parse(btn.dataset.sizes);
            updatePackSizes(sizes);
        });
    });
}

// Handle Calculate Form Submit
async function handleCalculate(e) {
    e.preventDefault();
    const quantity = parseInt(orderQuantityInput.value);
    if (isNaN(quantity) || quantity < 0) {
        showError('Please enter a valid non-negative number');
        return;
    }
    await calculatePacks(quantity);
}

// Calculate Packs
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

// Display Result
function displayResult(data, orderQuantity) {
    orderDisplay.textContent = orderQuantity.toLocaleString();

    let totalItems = 0;
    let totalPacks = 0;

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

// Load Pack Sizes
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
        // Use default sizes if API fails
        renderPackSizeInputs([250, 500, 1000, 2000, 5000]);
    }
}

// Render Pack Size Inputs
function renderPackSizeInputs(sizes) {
    packSizesContainer.innerHTML = '';
    sizes.forEach(size => addPackSizeInput(size));
}

// Add Pack Size Input
function addPackSizeInput(value = '') {
    const div = document.createElement('div');
    div.className = 'pack-size-input';
    div.innerHTML = `
        <input type="number" min="1" value="${value}" placeholder="Size" required>
        <span>items</span>
        <button type="button" class="btn btn-danger remove-pack" aria-label="Remove pack size">×</button>
    `;

    // Remove button handler
    div.querySelector('.remove-pack').addEventListener('click', () => {
        if (packSizesContainer.children.length > 1) {
            div.remove();
        } else {
            showError('At least one pack size is required');
        }
    });

    packSizesContainer.appendChild(div);
}

// Handle Update Config
async function handleUpdateConfig(e) {
    e.preventDefault();

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

// Update Pack Sizes
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

        // If there's an order quantity, recalculate with new sizes
        const currentQuantity = parseInt(orderQuantityInput.value);
        if (currentQuantity > 0) {
            await calculatePacks(currentQuantity);
        }
    } catch (error) {
        showError(error.message || 'An error occurred while updating pack sizes');
    }
}

// UI Helpers
function showResult() {
    resultSection.classList.remove('hidden');
}

function hideResult() {
    resultSection.classList.add('hidden');
}

function showError(message) {
    errorMessage.textContent = message;
    errorSection.classList.remove('hidden');
}

function hideError() {
    errorSection.classList.add('hidden');
}

// Handle input validation - prevent negative numbers
orderQuantityInput.addEventListener('input', function() {
    if (this.value < 0) {
        this.value = 0;
    }
});

// Add event delegation for dynamic pack size inputs
packSizesContainer.addEventListener('input', function(e) {
    if (e.target.type === 'number') {
        if (e.target.value < 1) {
            e.target.value = 1;
        }
    }
});
