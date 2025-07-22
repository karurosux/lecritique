<script>
  import { page } from "$app/stores";
  import {
    FeatureGate,
    FEATURES,
    ActiveSubscriptionGate,
  } from "$lib/components/subscription";
  import { Logo, UserMenu } from "$lib/components/ui";
  import { auth } from "$lib/stores/auth";
  import { subscription } from "$lib/stores/subscription";
  import "../app.css";

  let { children } = $props();

  let isAuthPage = $derived(
    $page.route?.id?.includes("login") ||
      $page.route?.id?.includes("register") ||
      $page.route?.id?.includes("forgot-password") ||
      $page.route?.id?.includes("reset-password") ||
      $page.route?.id?.includes("registration-success") ||
      $page.route?.id?.includes("email-verification") ||
      $page.route?.id?.includes("verify-email") ||
      $page.route?.id?.includes("team/accept-invite"),
  );

  let isPublicPage = $derived(
    $page.route?.id?.includes("qr/") ||
      $page.route?.id?.includes("feedback/success"),
  );

  let isLegalPage = $derived(
    $page.route?.id?.includes("terms") || $page.route?.id?.includes("privacy"),
  );

  // Show navbar on all non-auth, non-public, and non-legal pages (independent of auth loading state)
  let showNavbar = $derived(!isAuthPage && !isPublicPage && !isLegalPage);

  let authState = $derived($auth);

  // Define routes where animated background should appear
  let routesWithAnimatedBg = [
    "/", // Landing page
    "/pricing", // Pricing page
    // Add more routes as needed
  ];

  let showAnimatedBackground = $derived(
    routesWithAnimatedBg.some((route) => {
      if (route === "/") {
        return $page.route?.id === "/";
      }
      return $page.route?.id?.includes(route);
    }),
  );
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
            <ActiveSubscriptionGate>
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
            </ActiveSubscriptionGate>
            <ActiveSubscriptionGate>
              <a
                href="/organizations"
                class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-150 {$page.route?.id?.includes(
                  'organizations',
                )
                  ? 'bg-gradient-to-r from-blue-600 to-purple-600 !text-white shadow-lg shadow-blue-500/25'
                  : 'text-gray-700 hover:text-gray-900 hover:bg-gray-100'}"
              >
                Organizations
              </a>
            </ActiveSubscriptionGate>
            <ActiveSubscriptionGate>
              <FeatureGate feature={FEATURES.BASIC_ANALYTICS}>
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
              </FeatureGate>
            </ActiveSubscriptionGate>
            <ActiveSubscriptionGate>
              <FeatureGate feature={FEATURES.FEEDBACK_EXPLORER}>
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
              </FeatureGate>
            </ActiveSubscriptionGate>
          </nav>
          <UserMenu />
        </div>
      </div>
    </div>
  </div>
{/if}

<div class="min-h-screen relative">
  <!-- Animated gradient orbs with light colors -->
  {#if showAnimatedBackground}
    <div class="animated-orbs-container fixed inset-0 overflow-hidden">
      <div
        class="animated-orb orb-purple-pink absolute -top-[40%] -left-[20%] w-[60%] h-[60%] bg-gradient-to-br from-blue-600 to-purple-600 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-pulse"
      ></div>
      <div
        class="animated-orb orb-blue-cyan absolute -bottom-[40%] -right-[20%] w-[60%] h-[60%] bg-gradient-to-br from-blue-500 to-blue-600 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-pulse animation-delay-2000"
      ></div>
      <div
        class="animated-orb orb-indigo-purple absolute top-[20%] right-[30%] w-[40%] h-[40%] bg-gradient-to-br from-purple-600 to-purple-700 rounded-full mix-blend-multiply filter blur-3xl opacity-15 animate-pulse animation-delay-4000"
      ></div>
    </div>
  {/if}

  <!-- Background gradient matching navbar -->
  <div
    class="fixed inset-0 bg-gradient-to-r from-blue-500/5 to-purple-500/5"
  ></div>
  <div
    class="fixed inset-0 bg-gradient-to-br from-blue-500/3 to-purple-500/3"
  ></div>

  <!-- Subtle grid pattern -->
  <div class="global-grid-pattern fixed inset-0 opacity-[0.03] z-[1]">
    <div
      class="grid-lines absolute inset-0"
      style="background-image: repeating-linear-gradient(0deg, transparent, transparent 59px, rgba(0,0,0,1) 59px, rgba(0,0,0,1) 60px), repeating-linear-gradient(90deg, transparent, transparent 59px, rgba(0,0,0,1) 59px, rgba(0,0,0,1) 60px);"
    ></div>
  </div>

  <!-- Content -->
  <div class="relative z-10">
    {@render children()}
  </div>
</div>
