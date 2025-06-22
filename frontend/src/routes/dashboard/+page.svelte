<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { Button, Card } from '$lib/components/ui';
  import { goto } from '$app/navigation';

  $: authState = $auth;

  function handleLogout() {
    auth.logout();
    goto('/login');
  }
</script>

<svelte:head>
  <title>Dashboard - LeCritique</title>
  <meta name="description" content="LeCritique restaurant management dashboard" />
</svelte:head>

{#if authState.isAuthenticated && authState.user}
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">LeCritique Dashboard</h1>
            <p class="text-sm text-gray-600">Welcome back, {authState.user.company_name}!</p>
          </div>
          <Button variant="outline" on:click={handleLogout}>
            Logout
          </Button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <!-- Quick Stats -->
          <Card>
            <div class="text-center">
              <h3 class="text-lg font-medium text-gray-900 mb-2">Restaurants</h3>
              <p class="text-3xl font-bold text-blue-600">0</p>
              <p class="text-sm text-gray-500 mt-1">Manage your restaurants</p>
            </div>
          </Card>

          <Card>
            <div class="text-center">
              <h3 class="text-lg font-medium text-gray-900 mb-2">QR Codes</h3>
              <p class="text-3xl font-bold text-green-600">0</p>
              <p class="text-sm text-gray-500 mt-1">Active QR codes</p>
            </div>
          </Card>

          <Card>
            <div class="text-center">
              <h3 class="text-lg font-medium text-gray-900 mb-2">Feedback</h3>
              <p class="text-3xl font-bold text-purple-600">0</p>
              <p class="text-sm text-gray-500 mt-1">Total feedback received</p>
            </div>
          </Card>
        </div>

        <!-- Quick Actions -->
        <div class="mt-8">
          <h2 class="text-lg font-medium text-gray-900 mb-4">Quick Actions</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card>
              <div class="text-center py-6">
                <svg class="mx-auto h-12 w-12 text-blue-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m0 0H5m5 0v-4a1 1 0 011-1h2a1 1 0 011 1v4M7 7h3m3 0h3m-6 4h3m3 0h3"></path>
                </svg>
                <h3 class="text-lg font-medium text-gray-900 mb-2">Add Restaurant</h3>
                <p class="text-sm text-gray-500 mb-4">Create a new restaurant profile</p>
                <Button variant="primary" href="/restaurants/new">
                  Get Started
                </Button>
              </div>
            </Card>

            <Card>
              <div class="text-center py-6">
                <svg class="mx-auto h-12 w-12 text-green-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"></path>
                </svg>
                <h3 class="text-lg font-medium text-gray-900 mb-2">Generate QR Code</h3>
                <p class="text-sm text-gray-500 mb-4">Create QR codes for feedback collection</p>
                <Button variant="primary" href="/qr-codes/new">
                  Create QR Code
                </Button>
              </div>
            </Card>
          </div>
        </div>

        <!-- Account Status -->
        {#if !authState.user.email_verified}
          <div class="mt-8">
            <div class="bg-yellow-50 border border-yellow-200 rounded-md p-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-5 w-5 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                </div>
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-yellow-800">
                    Email Verification Required
                  </h3>
                  <div class="mt-1 text-sm text-yellow-700">
                    <p>Please verify your email address to access all features.</p>
                  </div>
                  <div class="mt-4">
                    <Button variant="outline" size="sm">
                      Resend Verification Email
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </main>
  </div>
{:else}
  <div class="min-h-screen bg-gray-50 flex items-center justify-center">
    <div class="text-center">
      <svg class="animate-spin h-8 w-8 text-blue-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <p class="text-gray-600">Loading...</p>
    </div>
  </div>
{/if}