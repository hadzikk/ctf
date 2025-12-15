// App initialization and core functionality

export function initializeApp() {
    // Initialize UI components
    setupTabNavigation();
    setupEventListeners();
    
    // Initialize any other app-wide functionality
    console.log('GIS CTF Application Initialized');
}

function setupTabNavigation() {
    const tabButtons = document.querySelectorAll('.tab-btn');
    
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            // Remove active class from all buttons and content
            document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));
            
            // Add active class to clicked button
            button.classList.add('active');
            
            // Show corresponding content
            const tabId = button.getAttribute('data-tab');
            document.getElementById(`${tabId}-tab`).classList.add('active');
        });
    });
}

function setupEventListeners() {
    // Navigation buttons
    document.getElementById('prev-challenge').addEventListener('click', navigateToPreviousChallenge);
    document.getElementById('next-challenge').addEventListener('click', navigateToNextChallenge);
    
    // Flag submission
    document.getElementById('submit-flag').addEventListener('click', submitFlag);
    document.getElementById('flag-input').addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            submitFlag();
        }
    });
}

function navigateToPreviousChallenge() {
    // Implementation will be added in the navigation module
    console.log('Navigate to previous challenge');
}

function navigateToNextChallenge() {
    // Implementation will be added in the navigation module
    console.log('Navigate to next challenge');
}

function submitFlag() {
    const flagInput = document.getElementById('flag-input');
    const feedbackElement = document.getElementById('flag-feedback');
    const flag = flagInput.value.trim();
    
    if (!flag) {
        showFeedback('Please enter a flag', 'error');
        return;
    }
    
    // In a real app, this would be an API call to verify the flag
    // For now, we'll simulate a response
    const isCorrect = Math.random() > 0.5; // 50% chance of success for demo
    
    if (isCorrect) {
        showFeedback('Correct flag! Challenge completed!', 'success');
        // Update UI to show challenge as completed
        const currentChallenge = document.querySelector('.challenge-item.active');
        if (currentChallenge) {
            currentChallenge.classList.add('completed');
        }
        // Clear input
        flagInput.value = '';
    } else {
        showFeedback('Incorrect flag. Try again!', 'error');
    }
}

function showFeedback(message, type) {
    const feedbackElement = document.getElementById('flag-feedback');
    feedbackElement.textContent = message;
    feedbackElement.className = ''; // Clear all classes
    feedbackElement.classList.add(type);
    feedbackElement.style.display = 'block';
    
    // Hide feedback after 5 seconds
    setTimeout(() => {
        feedbackElement.style.display = 'none';
    }, 5000);
}

// Export functions that need to be available to other modules
export { showFeedback };
