<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { Button } from '$lib/components/ui';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  $: authState = $auth;

  onMount(() => {
    // Redirect authenticated users to dashboard
    if (authState.isAuthenticated) {
      goto('/dashboard');
    }
  });
</script>

<svelte:head>
  <title>LeCritique - Restaurant Feedback Management</title>
  <meta name="description" content="Collect and manage customer feedback for your restaurant with QR code technology" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
  <!-- Header -->
  <header class="bg-white shadow">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center py-6">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">LeCritique</h1>
        </div>
        <div class="flex space-x-4">
          <Button variant="outline" href="/login">
            Sign In
          </Button>
          <Button variant="primary" href="/register">
            Get Started
          </Button>
        </div>
      </div>
    </div>
  </header>

  <!-- Hero Section -->
  <main>
    <div class="max-w-7xl mx-auto py-16 px-4 sm:py-24 sm:px-6 lg:px-8">
      <div class="text-center">
        <h2 class="text-4xl font-bold text-gray-900 sm:text-5xl md:text-6xl">
          Restaurant Feedback
          <span class="text-blue-600">Made Simple</span>
        </h2>
        <p class="mt-5 max-w-md mx-auto text-base text-gray-500 sm:text-lg md:mt-8 md:text-xl md:max-w-3xl">
          Collect customer feedback effortlessly with QR codes. Improve your restaurant's service and menu based on real customer insights.
        </p>
        <div class="mt-8 flex justify-center">
          <div class="inline-flex rounded-md shadow">
            <Button variant="primary" size="lg" href="/register">
              Start Free Trial
            </Button>
          </div>
          <div class="ml-3 inline-flex">
            <Button variant="outline" size="lg" href="/demo">
              View Demo
            </Button>
          </div>
        </div>
      </div>
    </div>

    <!-- Features Section -->
    <div class="py-16 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="lg:text-center">
          <h2 class="text-base text-blue-600 font-semibold tracking-wide uppercase">Features</h2>
          <p class="mt-2 text-3xl leading-8 font-bold tracking-tight text-gray-900 sm:text-4xl">
            Everything you need to manage feedback
          </p>
        </div>

        <div class="mt-10">
          <div class="space-y-10 md:space-y-0 md:grid md:grid-cols-2 md:gap-x-8 md:gap-y-10">
            <div class="relative">
              <div class="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"></path>
                </svg>
              </div>
              <p class="ml-16 text-lg leading-6 font-medium text-gray-900">QR Code Integration</p>
              <p class="mt-2 ml-16 text-base text-gray-500">
                Generate QR codes for tables, locations, or dishes. Customers scan and provide instant feedback.
              </p>
            </div>

            <div class="relative">
              <div class="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
                </svg>
              </div>
              <p class="ml-16 text-lg leading-6 font-medium text-gray-900">Analytics Dashboard</p>
              <p class="mt-2 ml-16 text-base text-gray-500">
                Track ratings, analyze trends, and get insights to improve your restaurant's performance.
              </p>
            </div>

            <div class="relative">
              <div class="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4"></path>
                </svg>
              </div>
              <p class="ml-16 text-lg leading-6 font-medium text-gray-900">Custom Questionnaires</p>
              <p class="mt-2 ml-16 text-base text-gray-500">
                Create tailored feedback forms for different dishes, services, or experiences.
              </p>
            </div>

            <div class="relative">
              <div class="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                </svg>
              </div>
              <p class="ml-16 text-lg leading-6 font-medium text-gray-900">Real-time Notifications</p>
              <p class="mt-2 ml-16 text-base text-gray-500">
                Get instantly notified about new feedback, especially critical reviews that need immediate attention.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </main>
</div>
