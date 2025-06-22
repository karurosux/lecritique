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

<div class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
  <div class="max-w-md w-full">
    <Card>
      <div class="text-center py-12 relative">
        <!-- Success Icon with Animation -->
        <div class="relative inline-block">
          <div class="absolute inset-0 bg-green-100 rounded-full animate-ping"></div>
          <svg 
            class="h-20 w-20 text-green-500 mx-auto mb-6 relative z-10" 
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              stroke-width="2" 
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" 
            />
          </svg>
        </div>

        <h1 class="text-3xl font-bold text-gray-900 mb-4">Thank You!</h1>
        
        <p class="text-lg text-gray-600 mb-6">
          Your feedback has been submitted successfully.
        </p>
        
        <div class="bg-blue-50 rounded-lg p-4 mb-6">
          <p class="text-blue-800 font-medium mb-2">
            Your opinion matters! 💙
          </p>
          <p class="text-sm text-blue-700">
            We appreciate you taking the time to share your experience. 
            Your feedback helps us improve our service and create better dining experiences.
          </p>
        </div>

        <!-- Social Sharing -->
        <div class="border-t border-gray-200 pt-6">
          <p class="text-sm text-gray-600 mb-4">
            Enjoyed your experience? Share it with others!
          </p>
          <div class="flex justify-center space-x-3">
            <button
              class="inline-flex items-center px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors"
              on:click={() => window.open('https://facebook.com', '_blank')}
            >
              <svg class="h-5 w-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
                <path d="M18 2h-3a5 5 0 00-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 011-1h3z" />
              </svg>
              Share
            </button>
            <button
              class="inline-flex items-center px-4 py-2 bg-black text-white rounded-md hover:bg-gray-800 transition-colors"
              on:click={() => window.open('https://twitter.com', '_blank')}
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
            on:click={handleBackToHome}
            class="w-full"
          >
            Back to Home
          </Button>
        </div>
      </div>
    </Card>

    <!-- Fun Facts -->
    <div class="mt-6 text-center">
      <p class="text-sm text-gray-500">
        <span class="font-medium">Fun fact:</span> 
        Your feedback helps restaurants improve their service by up to 40%! 📈
      </p>
    </div>
  </div>

  <!-- Confetti Effect (CSS Animation) -->
  {#if showConfetti}
    <div class="fixed inset-0 pointer-events-none overflow-hidden">
      {#each Array(50) as _, i}
        <div 
          class="absolute w-2 h-2 bg-{['red', 'blue', 'green', 'yellow', 'purple'][i % 5]}-500 animate-fall"
          style="
            left: {Math.random() * 100}%;
            animation-delay: {Math.random() * 3}s;
            animation-duration: {3 + Math.random() * 2}s;
          "
        ></div>
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
</style>