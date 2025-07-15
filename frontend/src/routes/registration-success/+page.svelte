<script lang="ts">
  import { Logo } from '$lib/components/ui';
  import { CheckCircle, Mail, ArrowRight } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import { APP_CONFIG, createMailtoLink } from '$lib/constants/config';
  
  let showContent = $state(false);
  
  onMount(() => {
    setTimeout(() => {
      showContent = true;
    }, 100);
  });
</script>

<svelte:head>
  <title>Registration Successful - LeCritique</title>
  <meta name="description" content="Your LeCritique account has been created successfully" />
</svelte:head>

<div class="min-h-screen bg-gradient-to-b from-white to-gray-50/50 flex items-center justify-center px-4">
  <div class="registration-success-container max-w-md w-full">
    <div class="text-center">
      <!-- Logo -->
      <div class="flex justify-center mb-8">
        <Logo size="lg" />
      </div>
      
      <!-- Success Icon with Animation -->
      <div class="mb-8 success-icon-container" class:show={showContent}>
        <div class="relative inline-flex">
          <div class="absolute inset-0 bg-green-500/20 rounded-full blur-xl animate-pulse"></div>
          <div class="relative bg-gradient-to-br from-green-400 to-green-600 p-6 rounded-full">
            <CheckCircle class="w-12 h-12 text-white" />
          </div>
        </div>
      </div>
      
      <!-- Content -->
      <div class="space-y-6 content-container" class:show={showContent}>
        <div>
          <h1 class="text-3xl sm:text-4xl font-bold text-gray-900 mb-3">
            Registration Successful!
          </h1>
          <p class="text-lg text-gray-600">
            Welcome to LeCritique. Your account has been created successfully.
          </p>
        </div>
        
        <!-- Email Verification Notice -->
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-6 text-left">
          <div class="flex gap-4">
            <div class="flex-shrink-0">
              <Mail class="w-6 h-6 text-blue-600" />
            </div>
            <div class="space-y-2">
              <h3 class="font-semibold text-gray-900">Check Your Email</h3>
              <p class="text-sm text-gray-700">
                We've sent a verification email to your registered email address. 
                Please check your inbox and click the verification link to activate your account.
              </p>
              <p class="text-sm text-gray-600 italic">
                Didn't receive the email? Check your spam folder or wait a few minutes.
              </p>
            </div>
          </div>
        </div>
        
        <!-- Actions -->
        <div class="space-y-3 pt-4">
          <a 
            href="/login" 
            class="inline-flex items-center justify-center w-full px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg font-medium hover:shadow-lg transition-all duration-300 hover:scale-[1.02] group"
          >
            Continue to Login
            <ArrowRight class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
          </a>
          
          <button 
            onclick={() => window.location.href = createMailtoLink('support')}
            class="text-sm text-gray-600 hover:text-gray-800 transition-colors"
          >
            Need help? Contact support
          </button>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .success-icon-container {
    opacity: 0;
    transform: scale(0.5);
    transition: all 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  }
  
  .success-icon-container.show {
    opacity: 1;
    transform: scale(1);
  }
  
  .content-container {
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.6s ease-out;
    transition-delay: 0.3s;
  }
  
  .content-container.show {
    opacity: 1;
    transform: translateY(0);
  }
  
  @keyframes pulse {
    0%, 100% {
      opacity: 1;
    }
    50% {
      opacity: 0.5;
    }
  }
</style>