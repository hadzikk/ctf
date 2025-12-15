// Challenges module - handles loading and managing CTF challenges

// Map configuration for challenges
const mapConfigs = {
    'gps-spoofing': {
        center: [-6.2088, 106.8456], // Jakarta, Indonesia
        zoom: 13,
        markers: [
            {
                position: [-6.2088, 106.8456],
                popup: 'Monas - The National Monument',
                icon: 'flag',
                color: 'red'
            },
            {
                position: [-6.1754, 106.8272],
                popup: 'Istiqlal Mosque',
                icon: 'mosque',
                color: 'blue'
            },
            {
                position: [-6.1256, 106.6556],
                popup: 'Tangerang',
                icon: 'map-marker-alt',
                color: 'green'
            }
        ]
    },
    'geo-json-injection': {
        center: [-7.2575, 112.7521], // Surabaya, Indonesia
        zoom: 12,
        geoJson: {
            type: 'FeatureCollection',
            features: [
                {
                    type: 'Feature',
                    properties: {
                        name: 'Surabaya City',
                        population: 3000000
                    },
                    geometry: {
                        type: 'Polygon',
                        coordinates: [
                            [
                                [112.65, -7.20],
                                [112.90, -7.20],
                                [112.90, -7.35],
                                [112.65, -7.35],
                                [112.65, -7.20]
                            ]
                        ]
                    }
                }
            ]
        }
    },
    'map-tile-hijacking': {
        center: [-6.5971, 106.8060], // Bogor, Indonesia
        zoom: 15,
        tileLayer: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
        customTiles: [
            {
                url: 'https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png',
                name: 'OpenTopoMap',
                attribution: 'Map data: &copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors, <a href="http://viewfinderpanoramas.org">SRTM</a> | Map style: &copy; <a href="https://opentopomap.org">OpenTopoMap</a> (<a href="https://creativecommons.org/licenses/by-sa/3.0/">CC-BY-SA</a>)'
            }
        ]
    },
    'location-privacy-bypass': {
        center: [-8.4095, 115.1889], // Bali, Indonesia
        zoom: 10,
        layers: [
            {
                type: 'circle',
                center: [-8.4095, 115.1889],
                radius: 5000, // 5km radius
                color: 'red',
                fillColor: '#f03',
                fillOpacity: 0.2
            },
            {
                type: 'marker',
                position: [-8.4095, 115.1889],
                popup: 'Restricted Area',
                icon: 'ban',
                color: 'red'
            }
        ]
    },
    'geofence-escape': {
        center: [-6.9147, 107.6098], // Bandung, Indonesia
        zoom: 12,
        geofence: {
            bounds: [
                [-6.8, 107.5], // Southwest corner
                [-7.0, 107.7]  // Northeast corner
            ],
            color: '#ff7800',
            weight: 2,
            fillOpacity: 0.1
        },
        pointsOfInterest: [
            {
                position: [-6.9039, 107.6186],
                name: 'Gedung Sate',
                type: 'landmark'
            },
            {
                position: [-6.9175, 107.6191],
                name: 'Alun-Alun Bandung',
                type: 'park'
            }
        ]
    }
};

// Initialize map for a challenge
function initMap(challengeId, containerId) {
    const config = mapConfigs[challengeId];
    if (!config) return null;
    
    // Create map
    const map = L.map(containerId).setView(config.center, config.zoom);
    
    // Add default tile layer
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        maxZoom: 19
    }).addTo(map);
    
    // Add custom tile layers if specified
    if (config.customTiles) {
        const baseLayers = {
            'OpenStreetMap': L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            })
        };
        
        const overlayLayers = {};
        
        config.customTiles.forEach((tile, index) => {
            const layer = L.tileLayer(tile.url, {
                attribution: tile.attribution
            });
            overlayLayers[tile.name || `Layer ${index + 1}`] = layer;
        });
        
        // Add layer control
        L.control.layers(baseLayers, overlayLayers).addTo(map);
    }
    
    // Add markers if specified
    if (config.markers) {
        config.markers.forEach(marker => {
            const customIcon = L.divIcon({
                html: `<i class="fas fa-${marker.icon || 'map-marker-alt'} fa-2x" style="color: ${marker.color || '#e74c3c'}"></i>`,
                className: 'custom-div-icon',
                iconSize: [30, 30],
                iconAnchor: [15, 30],
                popupAnchor: [0, -30]
            });
            
            L.marker(marker.position, { icon: customIcon })
                .addTo(map)
                .bindPopup(marker.popup || 'Marker');
        });
    }
    
    // Add GeoJSON if specified
    if (config.geoJson) {
        L.geoJSON(config.geoJson, {
            style: function(feature) {
                return {
                    color: '#4a83ec',
                    weight: 2,
                    opacity: 1,
                    fillColor: '#b8d1f3',
                    fillOpacity: 0.5
                };
            },
            onEachFeature: function(feature, layer) {
                if (feature.properties && feature.properties.name) {
                    layer.bindPopup(`<strong>${feature.properties.name}</strong><br>Population: ${feature.properties.population || 'N/A'}`);
                }
            }
        }).addTo(map);
    }
    
    // Add geofence if specified
    if (config.geofence) {
        const bounds = L.latLngBounds(config.geofence.bounds);
        L.rectangle(bounds, {
            color: config.geofence.color || '#ff0000',
            weight: config.geofence.weight || 3,
            fillColor: config.geofence.fillColor || '#ff0000',
            fillOpacity: config.geofence.fillOpacity || 0.1
        }).addTo(map).bindPopup('Restricted Area');
        
        // Fit bounds if needed
        if (config.geofence.fitBounds) {
            map.fitBounds(bounds);
        }
    }
    
    // Add points of interest if specified
    if (config.pointsOfInterest) {
        config.pointsOfInterest.forEach(poi => {
            const icon = L.divIcon({
                html: `<i class="fas fa-${poi.type === 'landmark' ? 'landmark' : 'tree'} fa-2x" style="color: ${poi.type === 'landmark' ? '#e67e22' : '#27ae60'}"></i>`,
                className: 'custom-div-icon',
                iconSize: [30, 30],
                iconAnchor: [15, 30],
                popupAnchor: [0, -30]
            });
            
            L.marker(poi.position, { icon })
                .addTo(map)
                .bindPopup(`<strong>${poi.name}</strong><br>Type: ${poi.type}`);
        });
    }
    
    // Add layers if specified
    if (config.layers) {
        config.layers.forEach(layer => {
            switch (layer.type) {
                case 'circle':
                    L.circle(layer.center, {
                        radius: layer.radius || 1000,
                        color: layer.color || '#3388ff',
                        fillColor: layer.fillColor || '#3388ff',
                        fillOpacity: layer.fillOpacity || 0.2
                    }).addTo(map);
                    break;
                case 'marker':
                    const icon = L.divIcon({
                        html: `<i class="fas fa-${layer.icon || 'map-marker-alt'} fa-2x" style="color: ${layer.color || '#e74c3c'}"></i>`,
                        className: 'custom-div-icon',
                        iconSize: [30, 30],
                        iconAnchor: [15, 30],
                        popupAnchor: [0, -30]
                    });
                    
                    L.marker(layer.position, { icon })
                        .addTo(map)
                        .bindPopup(layer.popup || 'Marker');
                    break;
            }
        });
    }
    
    return map;
}

export function loadChallenges() {
    // Sample GIS-themed challenges
    const challenges = [
        {
            id: 'gps-spoofing',
            title: 'GPS Spoofing',
            points: 100,
            description: 'You need to find the hidden location by analyzing the GPS coordinates in the request.',
            instructions: 'Inspect the network requests to find the hidden flag in the GPS coordinates.',
            setup: '1. Open your browser\'s developer tools (F12)\n2. Go to the Network tab\n3. Look for requests containing location data',
            challenge: 'The application is sending location data to an API endpoint. Find the endpoint and analyze the response to get the flag.',
            techniques: '1. Use browser developer tools\n2. Look for API endpoints in the Network tab\n3. Check for hidden data in responses',
            flag: 'CTF{gps_sp00f1ng_1s_fun}'
        },
        {
            id: 'geo-json-injection',
            title: 'GeoJSON Injection',
            points: 200,
            description: 'Exploit a vulnerability in the GeoJSON parsing to reveal the flag.',
            instructions: 'The application processes GeoJSON data. Find a way to inject malicious input to reveal the flag.',
            setup: '1. Locate where the application processes GeoJSON data\n2. Prepare a crafted GeoJSON payload',
            challenge: 'The application is vulnerable to injection through the GeoJSON parser. Craft a payload that will reveal the flag.',
            techniques: '1. Study GeoJSON structure\n2. Try injecting special characters\n3. Look for server-side template injection',
            flag: 'CTF{ge0j50n_1nj3ct10n_ftw}'
        },
        {
            id: 'map-tile-hijacking',
            title: 'Map Tile Hijacking',
            points: 150,
            description: 'The application is loading map tiles from an insecure source. Intercept and modify the tiles to find the flag.',
            instructions: 'Intercept the map tile requests and modify them to reveal the hidden flag.',
            setup: '1. Use a proxy like Burp Suite\n2. Intercept the map tile requests\n3. Modify the responses',
            challenge: 'The application loads map tiles from an insecure source. Find a way to intercept and modify the tile requests to get the flag.',
            techniques: '1. Use a web proxy\n2. Intercept and modify requests\n3. Look for hidden data in map tiles',
            flag: 'CTF{m4p_t1l3_h1j4ck3d}'
        },
        {
            id: 'location-privacy-bypass',
            title: 'Location Privacy Bypass',
            points: 250,
            description: 'The application has a privacy feature that should hide certain locations. Find a way to bypass this protection.',
            instructions: 'The application is supposed to hide sensitive locations, but there\'s a way to bypass this protection.',
            setup: '1. Analyze how the application handles location privacy\n2. Look for client-side validation',
            challenge: 'Find a way to access the hidden locations that should be restricted.',
            techniques: '1. Check for client-side validation only\n2. Try modifying request parameters\n3. Look for API endpoints that might not be properly secured',
            flag: 'CTF{pr1v4cy_1s_4n_1llus10n}'
        },
        {
            id: 'geofence-escape',
            title: 'Geofence Escape',
            points: 300,
            description: 'The application has a geofence that restricts certain actions. Find a way to bypass this restriction.',
            instructions: 'The application uses geofencing to restrict access to certain features. Find a way to bypass this restriction.',
            setup: '1. Understand how the geofencing is implemented\n2. Look for client-side validation',
            challenge: 'The application uses client-side geofencing that can be bypassed. Find the flag by bypassing the geofence.',
            techniques: '1. Modify location data in the browser\n2. Use developer tools to override geolocation\n3. Look for API endpoints that don\'t validate location server-side',
            flag: 'CTF{g30f3nc3_byp4ss3d}'
        }
    ];

    // Render challenge list
    renderChallengeList(challenges);
    
    // Load the first challenge by default
    if (challenges.length > 0) {
        displayChallenge(challenges[0]);
    }
}

function renderChallengeList(challenges) {
    const challengeList = document.getElementById('challenge-list');
    
    challenges.forEach((challenge, index) => {
        const li = document.createElement('li');
        li.className = 'challenge-item';
        li.dataset.id = challenge.id;
        li.dataset.index = index;
        
        li.innerHTML = `
            <span>${index + 1}. ${challenge.title}</span>
            <span class="challenge-points">${challenge.points} pts</span>
        `;
        
        li.addEventListener('click', () => displayChallenge(challenge));
        challengeList.appendChild(li);
    });
}

export function displayChallenge(challenge) {
    if (!challenge) return;
    
    // Update active state in the sidebar
    document.querySelectorAll('.challenge-item').forEach(item => {
        item.classList.remove('active');
        if (item.dataset.id === challenge.id) {
            item.classList.add('active');
        }
    });
    
    // Clear any existing map
    const mapContainer = document.getElementById('challenge-map');
    if (mapContainer) {
        mapContainer.innerHTML = '';
        
        // Show loading state
        const loadingElement = document.querySelector('.map-loading');
        if (loadingElement) {
            loadingElement.style.display = 'flex';
        }
        
        // Initialize the map after a short delay to ensure the container is ready
        setTimeout(() => {
            // Initialize the map for this challenge
            const map = initMap(challenge.id, 'challenge-map');
            
            // Hide loading state
            if (loadingElement) {
                loadingElement.style.display = 'none';
            }
            
            // Store the map instance for later reference
            window.challengeMap = map;
        }, 100);
    }
    
    // Update challenge header
    document.getElementById('challenge-title').textContent = challenge.title;
    document.getElementById('challenge-points').textContent = `${challenge.points} pts`;
    
    // Update challenge content
    const challengeContent = document.getElementById('challenge-content');
    challengeContent.innerHTML = `
        <p>${challenge.description}</p>
        <div class="map-container">
            <div id="challenge-map"></div>
            <div class="map-loading">
                <i class="fas fa-map-marked-alt fa-spin"></i>
                <p>Loading map...</p>
            </div>
        </div>
        <div class="challenge-instructions">
            <h3>Instructions</h3>
            <p>${challenge.instructions}</p>
        </div>
    `;
    
    // Update tab content
    document.getElementById('setup-content').textContent = challenge.setup;
    document.getElementById('challenge-details').textContent = challenge.challenge;
    document.getElementById('techniques-content').innerHTML = challenge.techniques.replace(/\n/g, '<br>');
    
    // Reset flag input and feedback
    document.getElementById('flag-input').value = '';
    document.getElementById('flag-feedback').style.display = 'none';
    
    // Update navigation buttons
    updateNavigationButtons(challenge);
}

function updateNavigationButtons(challenge) {
    const challengeItems = Array.from(document.querySelectorAll('.challenge-item'));
    const currentIndex = challengeItems.findIndex(item => item.dataset.id === challenge.id);
    
    const prevButton = document.getElementById('prev-challenge');
    const nextButton = document.getElementById('next-challenge');
    
    // Update previous button
    if (currentIndex > 0) {
        prevButton.disabled = false;
        prevButton.onclick = () => {
            const prevChallenge = challengeItems[currentIndex - 1];
            const challengeId = prevChallenge.dataset.id;
            const challenge = getChallengeById(challengeId);
            displayChallenge(challenge);
        };
    } else {
        prevButton.disabled = true;
    }
    
    // Update next button
    if (currentIndex < challengeItems.length - 1) {
        nextButton.disabled = false;
        nextButton.onclick = () => {
            const nextChallenge = challengeItems[currentIndex + 1];
            const challengeId = nextChallenge.dataset.id;
            const challenge = getChallengeById(challengeId);
            displayChallenge(challenge);
        };
    } else {
        nextButton.disabled = true;
    }
}

function getChallengeById(challengeId) {
    // In a real app, this would fetch from an API
    // For now, we'll just return a mock challenge
    const challenges = [
        // Same challenges as above...
    ];
    
    return challenges.find(c => c.id === challengeId);
}

// Function to verify flag (would be an API call in a real app)
function verifyFlag(challengeId, flag) {
    const challenge = getChallengeById(challengeId);
    if (!challenge) return false;
    
    return flag === challenge.flag;
}

// Export functions that need to be available to other modules
export { verifyFlag };
