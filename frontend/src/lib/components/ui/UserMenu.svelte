<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  let authState = $derived($auth);
  let user = $derived(authState.user);

  let showMenu = $state(false);

  function toggleMenu() {
    showMenu = !showMenu;
  }

  function closeMenu() {
    showMenu = false;
  }

  async function handleLogout() {
    closeMenu();
    await auth.logout();
    // Small delay to ensure auth state is propagated
    await new Promise(resolve => setTimeout(resolve, 100));
    // Use replace to ensure we don't keep the current page in history
    await goto('/login', { replaceState: true });
  }

  function handleProfile() {
    // Navigate to profile page when implemented
    closeMenu();
  }

  function handleSettings() {
    // Navigate to settings page when implemented
    closeMenu();
  }

  // Generate avatar from email
  function getAvatarText(email: string): string {
    return email.charAt(0).toUpperCase();
  }

  function getAvatarColor(email: string): string {
    const colors = [
      'from-blue-500 to-purple-600',
      'from-green-500 to-emerald-600', 
      'from-yellow-500 to-orange-600',
      'from-red-500 to-pink-600',
      'from-indigo-500 to-blue-600',
      'from-purple-500 to-indigo-600'
    ];
    
    let hash = 0;
    for (let i = 0; i < email.length; i++) {
      hash = email.charCodeAt(i) + ((hash << 5) - hash);
    }
    
    return colors[Math.abs(hash) % colors.length];
  }

  // Close menu when clicking outside
  function handleClickOutside(event: MouseEvent) {
    const target = event.target as Element;
    if (!target.closest('.user-menu-container')) {
      closeMenu();
    }
  }
</script>

<svelte:window onclick={handleClickOutside} />

{#if user}
  <div class="relative user-menu-container z-[9999]">
    <!-- User Menu Button -->
    <button
      type="button"
      class="flex items-center space-x-3 p-2 rounded-xl hover:bg-white/10 transition-all duration-200 group cursor-pointer"
      onclick={toggleMenu}
      aria-expanded={showMenu}
      aria-haspopup="true"
    >
      <!-- Avatar -->
      <div class="h-10 w-10 bg-gradient-to-br {getAvatarColor(user.email)} rounded-xl flex items-center justify-center shadow-lg group-hover:scale-105 transition-transform duration-200">
        <span class="text-white font-semibold text-sm">
          {getAvatarText(user.email)}
        </span>
      </div>
      
      <!-- User Info -->
      <div class="hidden sm:block text-left">
        <p class="text-sm font-semibold text-gray-900 truncate max-w-32">
          {user.company_name || 'User'}
        </p>
        <p class="text-xs text-gray-600 truncate max-w-32">
          {user.email}
        </p>
      </div>
      
      <!-- Chevron -->
      <svg 
        class="h-4 w-4 text-gray-600 transition-transform duration-200 {showMenu ? 'rotate-180' : ''}"
        fill="none" 
        stroke="currentColor" 
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>

    <!-- Dropdown Menu -->
    {#if showMenu}
      <div class="absolute right-0 top-full mt-2 w-64 bg-white/95 backdrop-blur-xl rounded-2xl shadow-xl shadow-gray-900/10 border border-white/20 overflow-hidden z-[9999]">
        <div class="p-4 border-b border-gray-100/50">
          <div class="flex items-center space-x-3">
            <div class="h-12 w-12 bg-gradient-to-br {getAvatarColor(user.email)} rounded-xl flex items-center justify-center shadow-lg">
              <span class="text-white font-semibold">
                {getAvatarText(user.email)}
              </span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold text-gray-900 truncate">
                {user.company_name || 'User'}
              </p>
              <p class="text-xs text-gray-600 truncate">
                {user.email}
              </p>
              {#if user.email_verified}
                <div class="flex items-center space-x-1 mt-1">
                  <svg class="h-3 w-3 text-green-500" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span class="text-xs text-green-600 font-medium">Verified</span>
                </div>
              {:else}
                <div class="flex items-center space-x-1 mt-1">
                  <svg class="h-3 w-3 text-yellow-500" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
                  </svg>
                  <span class="text-xs text-yellow-600 font-medium">Unverified</span>
                </div>
              {/if}
            </div>
          </div>
        </div>

        <div class="py-2">
          <button
            type="button"
            class="w-full flex items-center px-4 py-3 text-sm text-gray-700 hover:bg-gray-50/80 transition-colors duration-150 cursor-pointer"
            onclick={handleProfile}
          >
            <svg class="h-5 w-5 mr-3 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            Profile
          </button>
          
          <button
            type="button"
            class="w-full flex items-center px-4 py-3 text-sm text-gray-700 hover:bg-gray-50/80 transition-colors duration-150 cursor-pointer"
            onclick={handleSettings}
          >
            <svg class="h-5 w-5 mr-3 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            Settings
          </button>
        </div>

        <div class="border-t border-gray-100/50">
          <button
            type="button"
            class="w-full flex items-center px-4 py-3 text-sm text-red-600 hover:bg-red-50/80 transition-colors duration-150 cursor-pointer"
            onclick={handleLogout}
          >
            <svg class="h-5 w-5 mr-3 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Sign out
          </button>
        </div>
      </div>
    {/if}
  </div>
{/if}