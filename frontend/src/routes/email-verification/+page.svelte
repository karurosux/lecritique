<script lang="ts">
  import { page } from '$app/stores';
  import { Logo } from '$lib/components/ui';
  import { Mail, AlertCircle, RefreshCw, ArrowLeft } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import { APP_CONFIG } from '$lib/constants/config';
  import { getApiClient } from '$lib/api/client';
  
  let email = $derived($page.url.searchParams.get('email') || '');
  let showContent = $state(false);
  let resending = $state(false);
  let resendSuccess = $state(false);
  let resendError = $state('');
  
  onMount(() => {
    setTimeout(() => {
      showContent = true;
    }, 100);
  });
  
  async function resendVerification() {
    if (resending || !email) return;
    
    resending = true;
    resendSuccess = false;
    resendError = '';
    
    try {
      const api = getApiClient();
      const response = await api.api.v1AuthResendVerificationCreate({ email });
      
      if (response.data.success) {
        resendSuccess = true;
        
        // Reset success message after 5 seconds
        setTimeout(() => {
          resendSuccess = false;
        }, 5000);
      }
    } catch (error: any) {
      resendError = error.response?.data?.error?.message || 'Failed to send verification email';
      
      // Reset error message after 5 seconds
      setTimeout(() => {
        resendError = '';
      }, 5000);
    } finally {
      resending = false;
    }
  }
</script>

<svelte:head>
  <title>Email Verification Required - Kyooar</title>
  <meta name="description" content="Please verify your email address to access Kyooar" />
</svelte:head>

<div class="min-h-screen bg-gradient-to-b from-white to-gray-50/50 flex items-center justify-center px-4">
  <div class="email-verification-container max-w-md w-full">
    <div class="text-center">
      <!-- Logo -->
      <div class="flex justify-center mb-8">
        <Logo size="lg" />
      </div>
      
      <!-- Alert Icon with Animation -->
      <div class="mb-8 alert-icon-container" class:show={showContent}>
        <div class="relative inline-flex">
          <div class="absolute inset-0 bg-orange-500/20 rounded-full blur-xl animate-pulse"></div>
          <div class="relative bg-gradient-to-br from-orange-400 to-orange-600 p-6 rounded-full">
            <Mail class="w-12 h-12 text-white" />
          </div>
        </div>
      </div>
      
      <!-- Content -->
      <div class="space-y-6 content-container" class:show={showContent}>
        <div>
          <h1 class="text-3xl sm:text-4xl font-bold text-gray-900 mb-3">
            Verify Your Email
          </h1>
          <p class="text-lg text-gray-600">
            Please verify your email address to access your Kyooar account
          </p>
        </div>
        
        <!-- Email Verification Notice -->
        <div class="bg-orange-50 border border-orange-200 rounded-lg p-6 text-left">
          <div class="flex gap-4">
            <div class="flex-shrink-0">
              <AlertCircle class="w-6 h-6 text-orange-600" />
            </div>
            <div class="space-y-2">
              <h3 class="font-semibold text-gray-900">Email Verification Required</h3>
              <p class="text-sm text-gray-700">
                We've sent a verification email to <strong class="font-medium">{email}</strong>
              </p>
              <p class="text-sm text-gray-600">
                Please check your inbox and click the verification link to activate your account. 
                The link will expire in 24 hours.
              </p>
              <div class="pt-2">
                <p class="text-sm text-gray-600 italic">
                  Can't find the email? Check your spam folder or click below to resend.
                </p>
              </div>
            </div>
          </div>
        </div>
        
        {#if resendSuccess}
          <div class="bg-green-50 border border-green-200 rounded-lg p-4">
            <p class="text-sm text-green-800 font-medium">
              Verification email has been resent successfully!
            </p>
          </div>
        {/if}
        
        {#if resendError}
          <div class="bg-red-50 border border-red-200 rounded-lg p-4">
            <p class="text-sm text-red-800 font-medium">
              {resendError}
            </p>
          </div>
        {/if}
        
        <!-- Actions -->
        <div class="space-y-3 pt-4">
          <button 
            onclick={resendVerification}
            disabled={resending}
            class="inline-flex items-center justify-center w-full px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg font-medium hover:shadow-lg transition-all duration-300 hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 group"
          >
            {#if resending}
              <RefreshCw class="w-5 h-5 mr-2 animate-spin" />
              Resending...
            {:else}
              <Mail class="w-5 h-5 mr-2 group-hover:scale-110 transition-transform" />
              Resend Verification Email
            {/if}
          </button>
          
          <a 
            href="/login" 
            class="inline-flex items-center justify-center w-full px-6 py-3 bg-white text-gray-700 rounded-lg font-medium border border-gray-300 hover:bg-gray-50 transition-all duration-200 group"
          >
            <ArrowLeft class="w-5 h-5 mr-2 group-hover:-translate-x-1 transition-transform" />
            Back to Login
          </a>
          
          <div class="pt-4 border-t border-gray-200">
            <p class="text-sm text-gray-600">
              Having trouble? 
              <a href={`mailto:${APP_CONFIG.emails.support}`} class="font-medium text-blue-600 hover:text-blue-700">
                Contact support
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .alert-icon-container {
    opacity: 0;
    transform: scale(0.5);
    transition: all 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  }
  
  .alert-icon-container.show {
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
</style>