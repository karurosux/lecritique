<script>
  import '../app.css';
  import { page } from '$app/stores';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { Button, UserMenu, Logo } from '$lib/components/ui';

  $: authState = $auth;
  $: isAuthPage = $page.route?.id?.includes('login') || $page.route?.id?.includes('register');
  $: showNavbar = authState.isAuthenticated && !isAuthPage;
</script>

{#if showNavbar}
  <!-- Shared Navigation Header -->
  <div class="relative bg-gradient-to-r from-white/95 to-gray-50/95 backdrop-blur-xl border-b border-white/20 shadow-lg shadow-gray-900/5 z-50">
    <div class="absolute inset-0 bg-gradient-to-r from-blue-500/5 to-purple-500/5"></div>
    <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center py-8">
        <div class="flex items-center">
          <Logo size="lg" />
        </div>
        <div class="flex items-center space-x-8">
          <!-- Navigation Menu -->
          <nav class="flex space-x-4">
            <a 
              href="/dashboard" 
              class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-150 {$page.route?.id?.includes('dashboard') 
                ? 'bg-gradient-to-r from-blue-600 to-purple-600 text-white shadow-lg shadow-blue-500/25' 
                : 'text-gray-700 hover:text-gray-900 hover:bg-gray-100'}"
            >
              Dashboard
            </a>
            <a 
              href="/restaurants" 
              class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-150 {$page.route?.id?.includes('restaurants') 
                ? 'bg-gradient-to-r from-blue-600 to-purple-600 text-white shadow-lg shadow-blue-500/25' 
                : 'text-gray-700 hover:text-gray-900 hover:bg-gray-100'}"
            >
              Restaurants
            </a>
            <a 
              href="/analytics" 
              class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-150 {$page.route?.id?.includes('analytics') 
                ? 'bg-gradient-to-r from-blue-600 to-purple-600 text-white shadow-lg shadow-blue-500/25' 
                : 'text-gray-700 hover:text-gray-900 hover:bg-gray-100'}"
            >
              Analytics
            </a>
            <a 
              href="/feedback/manage" 
              class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-150 {$page.route?.id?.includes('feedback') 
                ? 'bg-gradient-to-r from-blue-600 to-purple-600 text-white shadow-lg shadow-blue-500/25' 
                : 'text-gray-700 hover:text-gray-900 hover:bg-gray-100'}"
            >
              Feedback
            </a>
          </nav>
          
          <UserMenu />
        </div>
      </div>
    </div>
  </div>
{/if}

<div class="min-h-screen bg-gradient-to-br from-gray-100 via-gray-200/50 to-gray-300/70">
  <slot />
</div>