<script lang="ts">
  import { Logo } from '$lib/components/ui';
  import { privacyContent } from '$lib/content/privacy';
  import { CalendarDays, FileText, ArrowLeft } from 'lucide-svelte';
  import { onMount } from 'svelte';

  let scrollY = $state(0);

  onMount(() => {
    const handleScroll = () => {
      scrollY = window.scrollY;
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  });
</script>

<svelte:head>
  <title>Privacy Policy - Kyooar</title>
  <meta
    name="description"
    content="Kyooar Privacy Policy - Learn how we collect, use, and protect your data on our organization feedback management platform" />
</svelte:head>

<div class="min-h-screen bg-gradient-to-b from-white to-gray-50/50">
  <!-- Floating Header -->
  <div
    class="privacy-header fixed top-0 left-0 right-0 z-40 bg-white/80 backdrop-blur-xl border-b border-gray-100 transition-all duration-300"
    class:shadow-lg={scrollY > 50}>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
      <div class="flex items-center justify-between">
        <a href="/register" class="flex items-center gap-3 group">
          <ArrowLeft
            class="w-5 h-5 text-gray-400 group-hover:text-gray-600 transition-colors" />
          <Logo size="md" />
        </a>
        <div class="flex items-center gap-6 text-sm">
          <div class="hidden sm:flex items-center gap-2 text-gray-500">
            <CalendarDays class="w-4 h-4" />
            <span>Updated {privacyContent.lastUpdated}</span>
          </div>
          <div class="hidden sm:flex items-center gap-2 text-gray-500">
            <FileText class="w-4 h-4" />
            <span>v{privacyContent.version}</span>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="pt-24 pb-16">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Hero Section -->
      <div class="text-center mb-16">
        <h1
          class="text-5xl sm:text-6xl font-bold bg-gradient-to-r from-gray-900 via-blue-800 to-purple-800 bg-clip-text text-transparent mb-6">
          Privacy Policy
        </h1>
        <p class="text-xl text-gray-600 max-w-2xl mx-auto">
          Your privacy matters. Learn how Kyooar protects and manages your data
        </p>
      </div>

      <!-- Main Content -->
      <main class="max-w-4xl mx-auto">
        <div class="prose prose-lg max-w-none">
          {#each privacyContent.sections as section, index}
            <section id="section-{index}" class="mb-12 scroll-mt-32">
              <h2
                class="text-2xl sm:text-3xl font-bold text-gray-900 mb-4 leading-tight">
                {section.title}
              </h2>
              <div class="text-gray-600 leading-relaxed whitespace-pre-line">
                {section.content}
              </div>
            </section>
          {/each}
        </div>

        <!-- Footer -->
        <div class="mt-16 pt-12 border-t border-gray-200">
          <div
            class="bg-gradient-to-br from-blue-50 to-purple-50 rounded-xl p-8">
            <p class="text-center text-gray-700 mb-6">
              By using Kyooar, you acknowledge that you have read and understood
              our Privacy Policy.
            </p>
            <div class="flex flex-col sm:flex-row gap-4 justify-center">
              <a
                href="/register"
                class="inline-flex items-center justify-center px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg font-medium hover:shadow-lg transition-all duration-300 hover:scale-105">
                Create Your Account
              </a>
              <a
                href="/login"
                class="inline-flex items-center justify-center px-6 py-3 bg-white text-gray-700 rounded-lg font-medium border border-gray-300 hover:bg-gray-50 transition-all duration-200">
                Sign In
              </a>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</div>

<style>
  .privacy-header {
    transform: translateY(0);
  }

  .prose h2 {
    @apply text-2xl sm:text-3xl font-bold text-gray-900 mb-4 leading-tight;
  }

  .prose div {
    @apply text-gray-600 leading-relaxed;
  }
</style>
