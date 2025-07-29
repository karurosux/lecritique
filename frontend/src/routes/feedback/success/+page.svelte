<script lang="ts">
  import { onMount } from 'svelte';
  import { Check, Heart, Home } from 'lucide-svelte';
  import { goto } from '$app/navigation';
  import { Button, Card } from '$lib/components/ui';

  let showConfetti = true;
  let showButton = false;

  onMount(() => {
    // Show button after balloons finish (around 8 seconds)
    setTimeout(() => {
      showButton = true;
    }, 8000);
  });

  function handleGoHome() {
    goto('/');
  }
</script>

<svelte:head>
  <title>Thank You! - Kyooar</title>
  <meta name="description" content="Thank you for your feedback" />
</svelte:head>

<div
  class="min-h-screen bg-gradient-to-br from-indigo-50 via-purple-50 to-pink-50 flex items-center justify-center px-4 relative overflow-hidden">
  <!-- Animated Background Elements -->
  <div class="absolute inset-0 opacity-40">
    <div
      class="absolute top-1/4 left-1/4 w-72 h-72 bg-gradient-to-r from-purple-300 to-pink-300 rounded-full mix-blend-multiply filter blur-xl animate-blob">
    </div>
    <div
      class="absolute top-1/3 right-1/4 w-72 h-72 bg-gradient-to-r from-yellow-300 to-orange-300 rounded-full mix-blend-multiply filter blur-xl animate-blob animation-delay-2000">
    </div>
    <div
      class="absolute bottom-1/4 left-1/3 w-72 h-72 bg-gradient-to-r from-pink-300 to-purple-300 rounded-full mix-blend-multiply filter blur-xl animate-blob animation-delay-4000">
    </div>
  </div>

  <div class="max-w-lg w-full relative z-10">
    <Card
      variant="glass"
      class="rounded-[2rem] p-12 shadow-2xl shadow-purple-500/10 animate-slide-up">
      <div class="text-center relative">
        <!-- Modern Success Icon -->
        <div class="relative inline-block mb-8">
          <!-- Outer Ring Animation -->
          <div
            class="absolute inset-0 w-32 h-32 bg-gradient-to-r from-emerald-400 to-teal-500 rounded-full animate-pulse-slow opacity-20">
          </div>
          <div
            class="absolute inset-2 w-28 h-28 bg-gradient-to-r from-emerald-400 to-teal-500 rounded-full animate-ping opacity-20">
          </div>

          <!-- Main Success Icon -->
          <div
            class="relative w-32 h-32 bg-gradient-to-br from-emerald-400 via-teal-500 to-cyan-500 rounded-full flex items-center justify-center shadow-2xl shadow-emerald-500/40 transform hover:scale-105 transition-all duration-500 animate-bounce-once">
            <div
              class="w-28 h-28 bg-gradient-to-br from-emerald-500 to-teal-600 rounded-full flex items-center justify-center">
              <Check
                class="h-14 w-14 text-white animate-draw-check"
                strokeWidth="2.5" />
            </div>
          </div>
        </div>

        <!-- Modern Typography -->
        <div class="space-y-6 mb-10">
          <h1
            class="text-5xl font-black bg-gradient-to-r from-gray-900 via-purple-900 to-purple-800 bg-clip-text text-transparent leading-tight animate-fade-in-up">
            Amazing!
          </h1>

          <p
            class="text-xl text-gray-700 font-medium leading-relaxed animate-fade-in-up animation-delay-300">
            Your feedback has been received
          </p>
        </div>

        <!-- Modern Appreciation Card -->
        <div
          class="bg-gradient-to-r from-purple-500/10 via-pink-500/10 to-orange-500/10 backdrop-blur-sm rounded-3xl p-8 border border-white/20 shadow-xl shadow-purple-500/10 animate-fade-in-up animation-delay-600">
          <div class="flex items-center justify-center mb-4">
            <div
              class="w-16 h-16 bg-gradient-to-br from-purple-400 to-pink-500 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/30 animate-float">
              <Heart class="h-8 w-8 text-white" fill="currentColor" />
            </div>
          </div>
          <h2
            class="text-2xl font-bold bg-gradient-to-r from-purple-600 to-pink-600 bg-clip-text text-transparent mb-3">
            Thank you!
          </h2>
          <p class="text-gray-600 leading-relaxed">
            Your insights help us create unforgettable dining experiences. Every
            opinion matters in our journey to excellence.
          </p>
        </div>

        <!-- Home Button - Appears after animation -->
        {#if showButton}
          <div class="mt-8 animate-fade-in-up">
            <Button
              variant="gradient"
              size="lg"
              onclick={handleGoHome}
              class="shadow-lg shadow-purple-500/25 hover:shadow-xl hover:shadow-purple-500/40 transform hover:scale-105">
              <Home class="h-5 w-5 mr-3" />
              Continue exploring
            </Button>
          </div>
        {/if}
      </div>
    </Card>
  </div>

  <!-- Confetti Effect (CSS Animation) -->
  <div class="fixed inset-0 pointer-events-none overflow-hidden z-50">
    {#each Array(50) as _, i}
      <div
        class="absolute animate-fall opacity-0"
        style="
          left: {Math.random() * 100}%;
          top: -100px;
          animation-delay: {Math.random() * 3}s;
          animation-duration: {3 + Math.random() * 1.5}s;
        ">
        <div
          class="w-3 h-3 bg-gradient-to-br {[
            'from-red-400 to-pink-600',
            'from-blue-400 to-indigo-600',
            'from-green-400 to-emerald-600',
            'from-yellow-400 to-orange-600',
            'from-purple-400 to-indigo-600',
          ][i % 5]} rounded-full shadow-lg">
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  /* Confetti Animation */
  @keyframes fall {
    0% {
      transform: translateY(-50vh) rotate(0deg);
      opacity: 0;
    }
    8% {
      transform: translateY(-20px) rotate(45deg);
      opacity: 1;
    }
    18% {
      opacity: 1;
    }
    33% {
      opacity: 0.8;
    }
    48% {
      opacity: 0.4;
    }
    68% {
      opacity: 0.1;
    }
    100% {
      transform: translateY(calc(100vh + 100px)) rotate(360deg);
      opacity: 0;
    }
  }

  .animate-fall {
    animation: fall linear forwards;
  }

  /* Modern Blob Animation */
  @keyframes blob {
    0% {
      transform: translate(0px, 0px) scale(1);
    }
    33% {
      transform: translate(30px, -50px) scale(1.1);
    }
    66% {
      transform: translate(-20px, 20px) scale(0.9);
    }
    100% {
      transform: translate(0px, 0px) scale(1);
    }
  }

  .animate-blob {
    animation: blob 7s infinite;
  }

  .animation-delay-2000 {
    animation-delay: 2s;
  }

  .animation-delay-4000 {
    animation-delay: 4s;
  }

  /* Card Slide Up Animation */
  @keyframes slide-up {
    from {
      opacity: 0;
      transform: translateY(60px) scale(0.95);
    }
    to {
      opacity: 1;
      transform: translateY(0) scale(1);
    }
  }

  .animate-slide-up {
    animation: slide-up 0.8s ease-out forwards;
  }

  /* Bounce Once Animation */
  @keyframes bounce-once {
    0%,
    20%,
    50%,
    80%,
    100% {
      transform: translateY(0);
    }
    40% {
      transform: translateY(-20px);
    }
    60% {
      transform: translateY(-10px);
    }
  }

  .animate-bounce-once {
    animation: bounce-once 2s ease-in-out 0.5s;
  }

  /* Checkmark Drawing Animation */
  @keyframes draw-check {
    to {
      stroke-dashoffset: 0;
    }
  }

  .animate-draw-check {
    animation: draw-check 1s ease-in-out 1.2s forwards;
  }

  /* Fade In Up with Stagger */
  @keyframes fade-in-up {
    from {
      opacity: 0;
      transform: translateY(30px);
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

  .animation-delay-300 {
    animation-delay: 0.3s;
  }

  .animation-delay-600 {
    animation-delay: 0.6s;
  }

  /* Slow Pulse Animation */
  @keyframes pulse-slow {
    0%,
    100% {
      opacity: 0.2;
      transform: scale(1);
    }
    50% {
      opacity: 0.4;
      transform: scale(1.05);
    }
  }

  .animate-pulse-slow {
    animation: pulse-slow 3s ease-in-out infinite;
  }

  /* Float Animation */
  @keyframes float {
    0%,
    100% {
      transform: translateY(0px);
    }
    50% {
      transform: translateY(-10px);
    }
  }

  .animate-float {
    animation: float 3s ease-in-out infinite;
  }

  /* Glassmorphism Enhancement */
  .backdrop-blur-2xl {
    backdrop-filter: blur(40px);
  }
</style>
