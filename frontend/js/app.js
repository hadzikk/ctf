// Main application entry point
import { initializeApp } from './modules/app.js';
import { loadChallenges } from './modules/challenges.js';

// Initialize the application when the DOM is fully loaded
document.addEventListener('DOMContentLoaded', () => {
    initializeApp();
    loadChallenges();
});
