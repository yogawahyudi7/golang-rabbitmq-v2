// Configuration
const API_BASE_URL = 'http://localhost:8080/api/v1/monitoring';
const REFRESH_INTERVAL = 5000; // 5 seconds
const MAX_CHART_POINTS = 20;

// Global state
let currentFilter = 'all';
let ratesChart = null;
let chartData = {
    labels: [],
    publishRates: [],
    consumeRates: []
};

// Initialize dashboard
document.addEventListener('DOMContentLoaded', () => {
    initChart();
    fetchData();
    setInterval(fetchData, REFRESH_INTERVAL);
});

// Initialize Chart.js
function initChart() {
    const ctx = document.getElementById('ratesChart').getContext('2d');
    ratesChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: chartData.labels,
            datasets: [
                {
                    label: 'Publish Rate (msg/s)',
                    data: chartData.publishRates,
                    borderColor: 'rgb(59, 130, 246)',
                    backgroundColor: 'rgba(59, 130, 246, 0.1)',
                    tension: 0.4,
                    fill: true
                },
                {
                    label: 'Consume Rate (msg/s)',
                    data: chartData.consumeRates,
                    borderColor: 'rgb(147, 51, 234)',
                    backgroundColor: 'rgba(147, 51, 234, 0.1)',
                    tension: 0.4,
                    fill: true
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            aspectRatio: 3,
            plugins: {
                legend: {
                    display: true,
                    position: 'top'
                },
                tooltip: {
                    mode: 'index',
                    intersect: false
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: 'Messages per Second'
                    }
                },
                x: {
                    title: {
                        display: true,
                        text: 'Time'
                    }
                }
            },
            animation: {
                duration: 750
            }
        }
    });
}

// Update chart with new data
function updateChart(publishRate, consumeRate) {
    const now = new Date();
    const timeLabel = now.toLocaleTimeString();

    chartData.labels.push(timeLabel);
    chartData.publishRates.push(publishRate);
    chartData.consumeRates.push(consumeRate);

    // Keep only last MAX_CHART_POINTS
    if (chartData.labels.length > MAX_CHART_POINTS) {
        chartData.labels.shift();
        chartData.publishRates.shift();
        chartData.consumeRates.shift();
    }

    ratesChart.update('none');
}

// Fetch all data
async function fetchData() {
    try {
        await Promise.all([
            fetchMetrics(),
            fetchActivity()
        ]);
        updateLastUpdated();
    } catch (error) {
        console.error('Error fetching data:', error);
        showError('Failed to connect to API. Make sure the server is running.');
    }
}

// Fetch metrics
async function fetchMetrics() {
    try {
        const response = await fetch(`${API_BASE_URL}/metrics/detailed`);
        if (!response.ok) throw new Error('Failed to fetch metrics');

        const result = await response.json();
        const data = result.data;

        // Update connection status
        const statusEl = document.querySelector('.status');
        statusEl.textContent = 'Connected';
        statusEl.className = 'status connected';

        // Update uptime
        document.getElementById('uptime').textContent = formatUptime(data.uptime_seconds);

        // Update cards
        document.getElementById('published').textContent = data.messages.published;
        document.getElementById('consumed').textContent = data.messages.consumed;
        document.getElementById('failed').textContent = data.messages.failed;
        document.getElementById('success-rate').textContent =
            data.messages.success_rate ? data.messages.success_rate.toFixed(2) + '%' : '0%';

        document.getElementById('publish-rate').textContent =
            data.performance.publish_rate_per_sec.toFixed(2);
        document.getElementById('consume-rate').textContent =
            data.performance.consume_rate_per_sec.toFixed(2);

        const inProgress = data.activity.publishing_in_progress + data.activity.consuming_in_progress;
        document.getElementById('in-progress').textContent = inProgress;

        document.getElementById('avg-processing').textContent =
            data.performance.avg_processing_time_ms + 'ms';

        // Update chart
        updateChart(
            data.performance.publish_rate_per_sec,
            data.performance.consume_rate_per_sec
        );

    } catch (error) {
        const statusEl = document.querySelector('.status');
        statusEl.textContent = 'Disconnected';
        statusEl.className = 'status disconnected';
        throw error;
    }
}

// Fetch activity and events
async function fetchActivity() {
    try {
        const endpoint = currentFilter === 'all'
            ? `${API_BASE_URL}/activity?limit=20`
            : `${API_BASE_URL}/events/${currentFilter}?limit=20`;

        const response = await fetch(endpoint);
        if (!response.ok) throw new Error('Failed to fetch activity');

        const result = await response.json();
        const events = result.data.recent_events || result.data.events || [];

        updateEventsTable(events);
    } catch (error) {
        showError('Failed to load events');
        throw error;
    }
}

// Update events table
function updateEventsTable(events) {
    const tbody = document.getElementById('events-tbody');
    const loading = document.getElementById('events-loading');
    const content = document.getElementById('events-content');
    const errorDiv = document.getElementById('events-error');

    loading.style.display = 'none';
    errorDiv.style.display = 'none';
    content.style.display = 'block';

    if (events.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" style="text-align: center; padding: 40px; color: #666;">No events found</td></tr>';
        return;
    }

    tbody.innerHTML = events.map(event => {
        const time = new Date(event.timestamp).toLocaleTimeString();
        const duration = event.duration_ms ? `${(event.duration_ms / 1000000).toFixed(2)}ms` : '-';
        const statusBadge = event.success
            ? '<span class="badge success">Success</span>'
            : '<span class="badge error">Failed</span>';

        let typeBadge = '';
        if (event.type === 'publish') {
            typeBadge = '<span class="badge publish">Publish</span>';
        } else if (event.type === 'consume') {
            typeBadge = '<span class="badge consume">Consume</span>';
        } else if (event.type === 'failed') {
            typeBadge = '<span class="badge failed">Failed</span>';
        } else {
            typeBadge = `<span class="badge">${event.type}</span>`;
        }

        return `
            <tr>
                <td>${time}</td>
                <td>${typeBadge}</td>
                <td>${event.routing_key || '-'}</td>
                <td>${duration}</td>
                <td>${statusBadge}</td>
            </tr>
        `;
    }).join('');
}

// Filter events
function filterEvents(type) {
    currentFilter = type;

    // Update button states
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');

    // Show loading
    const loading = document.getElementById('events-loading');
    const content = document.getElementById('events-content');
    loading.style.display = 'block';
    content.style.display = 'none';

    // Fetch filtered data
    fetchActivity();
}

// Show error message
function showError(message) {
    const errorDiv = document.getElementById('events-error');
    const loading = document.getElementById('events-loading');
    const content = document.getElementById('events-content');

    loading.style.display = 'none';
    content.style.display = 'none';
    errorDiv.style.display = 'block';
    errorDiv.textContent = message;
}

// Update last updated timestamp
function updateLastUpdated() {
    const now = new Date();
    document.getElementById('last-updated').textContent = now.toLocaleTimeString();
}

// Format uptime
function formatUptime(seconds) {
    if (!seconds || seconds < 0) return '-';

    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = Math.floor(seconds % 60);

    if (hours > 0) {
        return `${hours}h ${minutes}m ${secs}s`;
    } else if (minutes > 0) {
        return `${minutes}m ${secs}s`;
    } else {
        return `${secs}s`;
    }
}

// Global filter function
window.filterEvents = filterEvents;
