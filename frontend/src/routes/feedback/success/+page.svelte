<script lang="ts">
  import { Card, Button } from '$lib/components/ui';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let showConfetti = true;

  onMount(() => {
    // Hide confetti after 3 seconds
    setTimeout(() => {
      showConfetti = false;
    }, 3000);
  });

  function handleBackToHome() {
    goto('/');
  }
</script>

<svelte:head>
  <title>Thank You! - LeCritique</title>
  <meta name="description" content="Thank you for your feedback" />
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-indigo-50/50 flex items-center justify-center px-4">
  <div class="max-w-md w-full">
    <Card variant="elevated" hover interactive class="group transform transition-all duration-300 animate-fade-in-up">
      <div class="text-center py-12 relative">
        <!-- Success Icon with Animation -->
        <div class="relative inline-block mb-8">
          <div class="absolute inset-0 bg-green-100 rounded-full animate-ping"></div>
          <div class="h-24 w-24 bg-gradient-to-br from-green-400 to-emerald-600 rounded-full flex items-center justify-center shadow-2xl shadow-green-500/50 relative z-10 group-hover:scale-110 transition-transform duration-300">
            <svg 
              class="h-12 w-12 text-white" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path 
                stroke-linecap="round" 
                stroke-linejoin="round" 
                stroke-width="3" 
                d="M5 13l4 4L19 7" 
              />
            </svg>
          </div>
        </div>

        <h1 class="text-4xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent mb-4">
          Thank You!
        </h1>
        
        <p class="text-lg text-gray-600 mb-8">
          Your feedback has been submitted successfully.
        </p>
        
        <div class="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-2xl p-6 mb-8 border border-blue-100">
          <div class="flex items-center justify-center mb-3">
            <span class="text-4xl">💙</span>
          </div>
          <p class="text-blue-900 font-semibold mb-2 text-lg">
            Your opinion matters!
          </p>
          <p class="text-sm text-blue-700">
            We appreciate you taking the time to share your experience. 
            Your feedback helps us improve our service and create better dining experiences.
          </p>
        </div>

        <!-- Social Sharing -->
        <div class="border-t border-gray-100 pt-8">
          <p class="text-sm text-gray-600 mb-5">
            Enjoyed your experience? Share it with others!
          </p>
          <div class="flex justify-center gap-3">
            <button
              class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-xl hover:from-blue-600 hover:to-blue-700 transition-all duration-200 shadow-lg shadow-blue-500/25 hover:shadow-xl hover:shadow-blue-500/40 transform hover:scale-105"
              onclick={() => window.open('https://facebook.com', '_blank')}
            >
              <svg class="h-5 w-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
                <path d="M18 2h-3a5 5 0 00-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 011-1h3z" />
              </svg>
              Share
            </button>
            <button
              class="inline-flex items-center px-6 py-3 bg-gradient-to-r from-gray-800 to-black text-white rounded-xl hover:from-gray-900 hover:to-black transition-all duration-200 shadow-lg shadow-gray-800/25 hover:shadow-xl hover:shadow-gray-800/40 transform hover:scale-105"
              onclick={() => window.open('https://twitter.com', '_blank')}
            >
              <svg class="h-5 w-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
                <path d="M23 3a10.9 10.9 0 01-3.14 1.53 4.48 4.48 0 00-7.86 3v1A10.66 10.66 0 013 4s-4 9 5 13a11.64 11.64 0 01-7 2c9 5 20 0 20-11.5a4.5 4.5 0 00-.08-.83A7.72 7.72 0 0023 3z" />
              </svg>
              Tweet
            </button>
          </div>
        </div>

        <!-- Action Button -->
        <div class="mt-8">
          <Button 
            variant="outline" 
            size="lg"
            onclick={handleBackToHome}
            class="w-full"
          >
            <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a2 2 0 002 2h10a2 2 0 002-2V10" />
            </svg>
            Back to Home
          </Button>
        </div>
      </div>
    </Card>

    <!-- Fun Facts -->
    <div class="mt-8 text-center">
      <div class="inline-flex items-center px-4 py-2 rounded-full bg-white/80 backdrop-blur-sm border border-gray-200/50 shadow-sm">
        <svg class="h-5 w-5 mr-2 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-sm text-gray-700">
          <span class="font-semibold">Fun fact:</span> 
          Your feedback helps restaurants improve their service by up to 40%! 📈
        </p>
      </div>
    </div>
  </div>

  <!-- Confetti Effect (CSS Animation) -->
  {#if showConfetti}
    <div class="fixed inset-0 pointer-events-none overflow-hidden">
      {#each Array(50) as _, i}
        <div 
          class="absolute animate-fall"
          style="
            left: {Math.random() * 100}%;
            animation-delay: {Math.random() * 3}s;
            animation-duration: {3 + Math.random() * 2}s;
          "
        >
          <div class="w-3 h-3 bg-gradient-to-br {['from-red-400 to-pink-600', 'from-blue-400 to-indigo-600', 'from-green-400 to-emerald-600', 'from-yellow-400 to-orange-600', 'from-purple-400 to-indigo-600'][i % 5]} rounded-full shadow-lg"></div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  @keyframes fall {
    to {
      transform: translateY(100vh) rotate(360deg);
      opacity: 0;
    }
  }
  
  .animate-fall {
    animation: fall linear infinite;
  }

  @keyframes fade-in-up {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-fade-in-up {
    animation: fade-in-up 0.8s ease-out forwards;
    opacity: 0;
  }
</style>