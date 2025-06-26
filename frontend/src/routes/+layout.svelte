<script>
  import { page } from "$app/stores";
  import { Logo, UserMenu } from "$lib/components/ui";
  import { auth } from "$lib/stores/auth";
  import "../app.css";

  let isAuthPage = $derived(
    $page.route?.id?.includes("login") || $page.route?.id?.includes("register"),
  );

  // Show navbar on all non-auth pages (independent of auth loading state)
  let showNavbar = $derived(!isAuthPage);
</script>

{#if showNavbar}
  <!-- Shared Navigation Header -->
  <div
    class="relative bg-gradient-to-r from-white/95 to-gray-50/95 backdrop-blur-xl border-b border-white/20 shadow-lg shadow-gray-900/5 z-50"
  >
    <div
      class="absolute inset-0 bg-gradient-to-r from-blue-500/5 to-purple-500/5"
    ></div>
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
              class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-150 {$page.route?.id?.includes(
                'dashboard',
              )
                ? 'bg-gradient-to-r from-blue-600 to-purple-600 !text-white shadow-lg shadow-blue-500/25'
                : 'text-gray-700 hover:text-gray-900 hover:bg-gray-100'}"
            >
              Dashboard
            </a>
            <a
              href="/restaurants"
              class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-150 {$page.route?.id?.includes(
                'restaurants',
              )
                ? 'bg-gradient-to-r from-blue-600 to-purple-600 !text-white shadow-lg shadow-blue-500/25'
                : 'text-gray-700 hover:text-gray-900 hover:bg-gray-100'}"
            >
              Restaurants
            </a>
            <a
              href="/analytics"
              class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-150 {$page.route?.id?.includes(
                'analytics',
              )
                ? 'bg-gradient-to-r from-blue-600 to-purple-600 !text-white shadow-lg shadow-blue-500/25'
                : 'text-gray-700 hover:text-gray-900 hover:bg-gray-100'}"
            >
              Analytics
            </a>
            <a
              href="/feedback/manage"
              class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-150 {$page.route?.id?.includes(
                'feedback',
              )
                ? 'bg-gradient-to-r from-blue-600 to-purple-600 !text-white shadow-lg shadow-blue-500/25'
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

<div
  class="min-h-screen bg-gradient-to-br from-blue-50/50 via-purple-50/30 to-pink-50/50 relative"
>
  <!-- Animated background gradients -->
  <div class="absolute inset-0 overflow-hidden">
    <div class="absolute -top-40 -right-40 w-80 h-80 bg-purple-300 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob"></div>
    <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-blue-300 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-2000"></div>
    <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-pink-300 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-4000"></div>
  </div>
  
  <!-- Content -->
  <div class="relative">
    <slot />
  </div>
</div>
