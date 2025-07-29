<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { Button } from '$lib/components/ui';
  import { goto } from '$app/navigation';
  import {
    QrCode,
    ArrowRight,
    Check,
    Building2,
    Users,
    Zap,
    BarChart3,
    Shield,
    Globe2,
  } from 'lucide-svelte';
  import { Logo } from '$lib/components/ui';
  import type { PageData } from './$types';
  import { PlanSelector } from '$lib/components/subscription';

  let { data }: { data: PageData } = $props();
  let authState = $derived($auth);

  $effect(() => {
    // Redirect authenticated users to dashboard
    if (authState.isAuthenticated) {
      goto('/dashboard');
    }
  });

  const features = [
    {
      title: 'Smart QR Codes',
      description: 'Generate QR codes that adapt to your business needs',
      icon: QrCode,
    },
    {
      title: 'Real-time Analytics',
      description:
        'Track customer interactions and engagement metrics instantly',
      icon: BarChart3,
    },
    {
      title: 'Multi-location Support',
      description: 'Manage multiple venues and locations from one dashboard',
      icon: Building2,
    },
    {
      title: 'Team Collaboration',
      description: 'Invite team members and manage permissions efficiently',
      icon: Users,
    },
    {
      title: 'Instant Updates',
      description: 'Changes reflect immediately across all QR codes',
      icon: Zap,
    },
    {
      title: 'Customer Feedback',
      description: 'Collect valuable insights directly from your customers',
      icon: Shield,
    },
  ];

  // Handle plan selection - redirect to registration with plan code
  function handleSelectPlan(plan: any) {
    goto(`/register?plan=${plan.code}`);
  }
</script>

<svelte:head>
  <title
    >Kyooar - Smart QR Code Platform for Customer Feedback & Analytics</title>
  <meta
    name="description"
    content="Transform customer feedback with smart QR codes. Collect valuable insights, analyze sentiment, and grow your business with data-driven decisions. Perfect for restaurants, retail, events & more." />
  <meta
    name="keywords"
    content="QR code platform, customer feedback, business analytics, restaurant feedback, retail analytics, customer insights, QR code generator, feedback management, customer satisfaction, business intelligence" />
  <meta name="author" content="Kyooar" />
  <meta name="robots" content="index, follow" />
  <link rel="canonical" href="https://kyooar.com" />

  <!-- Open Graph / Facebook -->
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://kyooar.com" />
  <meta
    property="og:title"
    content="Kyooar - Smart QR Code Platform for Customer Feedback & Analytics" />
  <meta
    property="og:description"
    content="Transform customer feedback with smart QR codes. Collect valuable insights, analyze sentiment, and grow your business with data-driven decisions." />
  <meta property="og:image" content="https://kyooar.com/og-image.jpg" />
  <meta property="og:site_name" content="Kyooar" />
  <meta property="og:locale" content="en_US" />

  <!-- Twitter -->
  <meta property="twitter:card" content="summary_large_image" />
  <meta property="twitter:url" content="https://kyooar.com" />
  <meta
    property="twitter:title"
    content="Kyooar - Smart QR Code Platform for Customer Feedback & Analytics" />
  <meta
    property="twitter:description"
    content="Transform customer feedback with smart QR codes. Collect valuable insights, analyze sentiment, and grow your business with data-driven decisions." />
  <meta
    property="twitter:image"
    content="https://kyooar.com/twitter-image.jpg" />
  <meta property="twitter:creator" content="@kyooar" />

  <!-- Additional SEO meta tags -->
  <meta name="theme-color" content="#3B82F6" />
  <meta name="msapplication-TileColor" content="#3B82F6" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />

  <!-- Structured Data -->
  <script type="application/ld+json">
    {
      "@context": "https://schema.org",
      "@type": "SoftwareApplication",
      "name": "Kyooar",
      "description": "Smart QR code platform for customer feedback and business analytics",
      "url": "https://kyooar.com",
      "applicationCategory": "BusinessApplication",
      "operatingSystem": "Web Browser",
      "offers": {
        "@type": "AggregateOffer",
        "priceCurrency": "USD",
        "lowPrice": "29.99",
        "highPrice": "199.99",
        "offerCount": "3"
      },
      "provider": {
        "@type": "Organization",
        "name": "Kyooar",
        "url": "https://kyooar.com",
        "logo": "https://kyooar.com/logo.png",
        "sameAs": [
          "https://twitter.com/kyooar",
          "https://linkedin.com/company/kyooar"
        ]
      },
      "featureList": [
        "Smart QR Code Generation",
        "Real-time Analytics",
        "Customer Feedback Collection",
        "Multi-location Support",
        "Team Collaboration",
        "Instant Updates"
      ]
    }
  </script>
</svelte:head>

<div class="landing-page min-h-screen flex flex-col">
  <!-- Navigation -->
  <nav
    class="fixed top-0 w-full bg-gradient-to-r from-white/95 to-gray-50/95 backdrop-blur-xl border-b border-white/20 shadow-lg shadow-gray-900/5 z-50">
    <div
      class="absolute inset-0 bg-gradient-to-r from-blue-500/5 to-purple-500/5">
    </div>
    <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-20">
        <Logo size="lg" />
        <div class="flex items-center space-x-8">
          <nav
            class="hidden md:flex items-center space-x-1"
            aria-label="Main navigation">
            <a
              href="#features"
              class="text-gray-700 hover:text-gray-900 font-medium px-4 py-2 rounded-xl hover:bg-white/50 transition-all duration-200 scroll-smooth"
              >Features</a>
            <a
              href="#use-cases"
              class="text-gray-700 hover:text-gray-900 font-medium px-4 py-2 rounded-xl hover:bg-white/50 transition-all duration-200 scroll-smooth"
              >Use Cases</a>
            <a
              href="#pricing"
              class="text-gray-700 hover:text-gray-900 font-medium px-4 py-2 rounded-xl hover:bg-white/50 transition-all duration-200 scroll-smooth"
              >Pricing</a>
            <a
              href="#contact"
              class="text-gray-700 hover:text-gray-900 font-medium px-4 py-2 rounded-xl hover:bg-white/50 transition-all duration-200 scroll-smooth"
              >Contact</a>
          </nav>
          <div class="flex items-center space-x-4">
            <Button variant="ghost" href="/login">Sign In</Button>
            <Button variant="gradient" href="/register">Get started</Button>
          </div>
        </div>
      </div>
    </div>
  </nav>

  <!-- Main Content -->
  <main class="flex-1">
    <!-- Hero Section -->
    <section
      class="relative min-h-screen flex items-center justify-center px-4 sm:px-6 lg:px-8 overflow-hidden"
      aria-label="Hero section">
      <!-- Background Elements -->
      <div class="absolute inset-0">
        <!-- Gradient mesh background -->
        <div
          class="absolute inset-0 bg-gradient-to-br from-blue-50 via-white to-purple-50">
        </div>

        <!-- Floating geometric shapes -->
        <div
          class="absolute top-20 left-10 w-20 h-20 bg-gradient-to-br from-blue-400 to-blue-600 rounded-2xl rotate-12 opacity-20 animate-float">
        </div>
        <div
          class="absolute top-40 right-20 w-16 h-16 bg-gradient-to-br from-purple-400 to-purple-600 rounded-full opacity-20 animate-float-delayed">
        </div>
        <div
          class="absolute bottom-40 left-20 w-12 h-12 bg-gradient-to-br from-pink-400 to-pink-600 rounded-xl -rotate-12 opacity-20 animate-float">
        </div>
        <div
          class="absolute bottom-20 right-40 w-24 h-24 bg-gradient-to-br from-indigo-400 to-indigo-600 rounded-3xl rotate-45 opacity-20 animate-float-delayed">
        </div>

        <!-- Additional animated elements -->
        <div
          class="absolute top-60 left-1/4 w-8 h-8 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-full opacity-15 animate-bounce-slow">
        </div>
        <div
          class="absolute top-32 right-1/3 w-14 h-14 bg-gradient-to-br from-yellow-400 to-orange-500 rounded-2xl rotate-45 opacity-15 animate-spin-slow">
        </div>
        <div
          class="absolute bottom-60 right-1/4 w-10 h-10 bg-gradient-to-br from-rose-400 to-rose-600 rounded-xl opacity-15 animate-pulse-slow">
        </div>
        <div
          class="absolute top-1/2 left-16 w-6 h-6 bg-gradient-to-br from-cyan-400 to-cyan-600 rounded-full opacity-15 animate-float-reverse">
        </div>
        <div
          class="absolute bottom-32 left-1/3 w-18 h-18 bg-gradient-to-br from-violet-400 to-violet-600 rounded-3xl rotate-12 opacity-15 animate-wiggle">
        </div>

        <!-- Floating particles -->
        <div
          class="absolute top-1/3 right-16 w-3 h-3 bg-blue-400 rounded-full opacity-30 animate-float-up">
        </div>
        <div
          class="absolute top-2/3 left-1/2 w-2 h-2 bg-purple-400 rounded-full opacity-25 animate-float-up-delayed">
        </div>
        <div
          class="absolute bottom-1/3 right-1/3 w-4 h-4 bg-pink-400 rounded-full opacity-20 animate-float-up">
        </div>
        <div
          class="absolute top-1/4 right-1/2 w-2 h-2 bg-indigo-400 rounded-full opacity-30 animate-float-up-slow">
        </div>

        <!-- Gradient orbs -->
        <div
          class="absolute top-1/4 left-1/2 w-32 h-32 bg-gradient-to-br from-blue-300/10 to-purple-300/10 rounded-full blur-xl animate-pulse-glow">
        </div>
        <div
          class="absolute bottom-1/4 right-1/2 w-40 h-40 bg-gradient-to-br from-pink-300/8 to-rose-300/8 rounded-full blur-xl animate-pulse-glow-delayed">
        </div>

        <!-- QR Code Scanning Animations -->
        <div class="absolute top-1/3 left-1/4 opacity-20">
          <div class="relative animate-qr-scan">
            <div
              class="w-16 h-16 border-2 border-blue-500 rounded-lg bg-white/80 backdrop-blur-sm p-2">
              <div
                class="w-full h-full bg-gradient-to-br from-blue-600 to-purple-600 rounded opacity-60">
              </div>
            </div>
            <div
              class="absolute inset-0 border-2 border-green-500 rounded-lg animate-scan-line">
            </div>
          </div>
        </div>

        <div class="absolute bottom-1/3 right-1/4 opacity-20">
          <div class="relative animate-qr-scan-delayed">
            <div
              class="w-12 h-12 border-2 border-purple-500 rounded-lg bg-white/80 backdrop-blur-sm p-1.5">
              <div
                class="w-full h-full bg-gradient-to-br from-purple-600 to-pink-600 rounded opacity-60">
              </div>
            </div>
            <div
              class="absolute inset-0 border-2 border-yellow-500 rounded-lg animate-scan-line-delayed">
            </div>
          </div>
        </div>

        <div class="absolute top-2/3 left-1/3 opacity-15">
          <div class="relative animate-qr-scan-slow">
            <div
              class="w-20 h-20 border-2 border-indigo-500 rounded-lg bg-white/80 backdrop-blur-sm p-2.5">
              <div
                class="w-full h-full bg-gradient-to-br from-indigo-600 to-blue-600 rounded opacity-60">
              </div>
            </div>
            <div
              class="absolute inset-0 border-2 border-emerald-500 rounded-lg animate-scan-pulse">
            </div>
          </div>
        </div>

        <div class="absolute top-1/2 right-1/3 opacity-10">
          <div class="relative animate-qr-scan-reverse">
            <div
              class="w-14 h-14 border-2 border-rose-500 rounded-lg bg-white/80 backdrop-blur-sm p-2">
              <div
                class="w-full h-full bg-gradient-to-br from-rose-600 to-red-600 rounded opacity-60">
              </div>
            </div>
            <div
              class="absolute inset-0 border-2 border-cyan-500 rounded-lg animate-scan-bounce">
            </div>
          </div>
        </div>

        <!-- Scanning laser effects -->
        <div
          class="absolute top-1/4 left-1/2 w-32 h-0.5 bg-gradient-to-r from-transparent via-green-400 to-transparent opacity-30 animate-laser-sweep">
        </div>
        <div
          class="absolute bottom-1/4 right-1/2 w-24 h-0.5 bg-gradient-to-r from-transparent via-blue-400 to-transparent opacity-25 animate-laser-sweep-delayed">
        </div>

        <!-- Grid pattern -->
        <div class="absolute inset-0 opacity-[0.02]">
          <div class="grid-pattern w-full h-full"></div>
        </div>
      </div>

      <div class="relative z-10 max-w-6xl mx-auto text-center pt-24">
        <!-- Main headline with animation -->
        <div class="mb-8">
          <h1
            class="text-4xl sm:text-5xl lg:text-6xl xl:text-7xl font-black leading-[0.9] tracking-tight">
            <span class="block text-gray-900 animate-slide-up">Transform</span>
            <span
              class="block bg-gradient-to-r from-blue-600 via-purple-600 to-pink-600 bg-clip-text text-transparent animate-slide-up-delayed"
              >Customer Feedback</span>
          </h1>
          <h2
            class="mt-4 text-xl sm:text-2xl font-bold text-gray-700 animate-fade-in-delayed">
            with Smart QR Codes
          </h2>
        </div>

        <!-- Description -->
        <p
          class="text-xl sm:text-2xl text-gray-600 mb-12 max-w-4xl mx-auto leading-relaxed font-light animate-fade-in-delayed">
          Collect valuable insights, analyze customer sentiment, and grow your
          business with data-driven decisions.
        </p>

        <!-- CTA Button -->
        <div class="flex justify-center mb-16 animate-fade-in-delayed">
          <Button
            variant="gradient"
            size="xl"
            href="/register"
            class="group min-w-[250px] transform hover:scale-105 transition-all duration-200 shadow-2xl hover:shadow-blue-500/25"
            aria-describedby="cta-description">
            Start Building Today
            <ArrowRight
              class="w-6 h-6 ml-3 group-hover:translate-x-2 transition-transform duration-200"
              aria-hidden="true" />
          </Button>
        </div>

        <!-- Social proof / Quick stats -->
        <div
          id="cta-description"
          class="flex flex-wrap justify-center gap-8 text-gray-500 text-sm animate-fade-in-delayed">
          <div class="flex items-center">
            <div
              class="w-2 h-2 bg-green-500 rounded-full mr-2"
              aria-hidden="true">
            </div>
            <span>No setup fees</span>
          </div>
          <div class="flex items-center">
            <div
              class="w-2 h-2 bg-green-500 rounded-full mr-2"
              aria-hidden="true">
            </div>
            <span>Ready in minutes</span>
          </div>
          <div class="flex items-center">
            <div
              class="w-2 h-2 bg-green-500 rounded-full mr-2"
              aria-hidden="true">
            </div>
            <span>Cancel anytime</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section
      id="features"
      class="py-32 relative overflow-hidden"
      aria-labelledby="features-heading">
      <div class="absolute inset-0" aria-hidden="true">
        <div
          class="absolute top-20 left-20 w-72 h-72 bg-blue-500/10 rounded-full blur-3xl">
        </div>
        <div
          class="absolute bottom-20 right-20 w-96 h-96 bg-purple-500/10 rounded-full blur-3xl">
        </div>
      </div>
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <header class="text-center mb-20">
          <div
            class="inline-flex items-center px-4 py-2 bg-gradient-to-r from-blue-500/10 to-purple-500/10 backdrop-blur-sm border border-blue-500/20 rounded-full mb-6">
            <span
              class="text-sm font-semibold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent"
              >Features that scale</span>
          </div>
          <h2
            id="features-heading"
            class="text-5xl sm:text-6xl font-black text-gray-900 mb-6 tracking-tight">
            Built for
            <span
              class="bg-gradient-to-r from-blue-600 via-purple-600 to-pink-600 bg-clip-text text-transparent"
              >modern business</span>
          </h2>
          <p class="text-xl text-gray-600 max-w-3xl mx-auto leading-relaxed">
            Everything you need to collect, analyze, and act on customer
            feedback at scale.
          </p>
        </header>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {#each features as feature, i}
            {@const gradients = [
              { bg: 'from-blue-500 to-purple-600' },
              { bg: 'from-purple-500 to-pink-600' },
              { bg: 'from-green-500 to-emerald-600' },
              { bg: 'from-yellow-500 to-orange-600' },
              { bg: 'from-indigo-500 to-blue-600' },
              { bg: 'from-red-500 to-pink-600' },
            ]}
            {@const gradient = gradients[i % gradients.length]}

            <div class="group relative">
              <!-- Glow effect -->
              <div
                class="absolute -inset-0.5 bg-gradient-to-r {gradient.bg} rounded-3xl blur opacity-0 group-hover:opacity-30 transition duration-1000 group-hover:duration-200">
              </div>

              <!-- Card -->
              <div
                class="relative bg-white/90 backdrop-blur-xl rounded-3xl p-8 border border-white/60 shadow-xl hover:shadow-2xl hover:-translate-y-2 transition-all duration-500">
                <!-- Icon -->
                <div class="relative mb-6">
                  <div
                    class="h-16 w-16 bg-gradient-to-r {gradient.bg} rounded-2xl flex items-center justify-center shadow-lg transform group-hover:scale-110 transition-transform duration-300">
                    <svelte:component
                      this={feature.icon}
                      class="w-8 h-8 text-white" />
                  </div>
                </div>

                <!-- Content -->
                <h3
                  class="text-xl font-bold text-gray-900 mb-3 group-hover:text-gray-800 transition-colors">
                  {feature.title}
                </h3>
                <p
                  class="text-gray-600 leading-relaxed group-hover:text-gray-700 transition-colors">
                  {feature.description}
                </p>

                <!-- Hover indicator -->
                <div
                  class="absolute bottom-6 right-6 w-2 h-2 bg-gradient-to-r {gradient.bg} rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </section>

    <!-- Use Cases -->
    <section id="use-cases" class="py-20" aria-labelledby="use-cases-heading">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
          <div>
            <h2
              id="use-cases-heading"
              class="text-3xl sm:text-4xl font-bold text-gray-900 mb-6">
              Built for every industry
            </h2>
            <p class="text-xl text-gray-600 mb-8">
              Whether you're running a restaurant, retail store, or enterprise
              organization, Kyooar adapts to your unique needs.
            </p>

            <ul class="space-y-4">
              <li class="flex items-start">
                <Check
                  class="w-6 h-6 text-green-500 mr-3 flex-shrink-0 mt-0.5"
                  aria-hidden="true" />
                <div>
                  <h3 class="font-semibold text-gray-900">
                    Restaurants & Cafes
                  </h3>
                  <p class="text-gray-600">
                    Digital menus, customer feedback, and table ordering
                  </p>
                </div>
              </li>

              <li class="flex items-start">
                <Check
                  class="w-6 h-6 text-green-500 mr-3 flex-shrink-0 mt-0.5"
                  aria-hidden="true" />
                <div>
                  <h3 class="font-semibold text-gray-900">
                    Retail & E-commerce
                  </h3>
                  <p class="text-gray-600">
                    Product information, reviews, and instant checkout
                  </p>
                </div>
              </li>

              <li class="flex items-start">
                <Check
                  class="w-6 h-6 text-green-500 mr-3 flex-shrink-0 mt-0.5"
                  aria-hidden="true" />
                <div>
                  <h3 class="font-semibold text-gray-900">
                    Events & Conferences
                  </h3>
                  <p class="text-gray-600">
                    Registration, networking, and session feedback
                  </p>
                </div>
              </li>

              <li class="flex items-start">
                <Check
                  class="w-6 h-6 text-green-500 mr-3 flex-shrink-0 mt-0.5"
                  aria-hidden="true" />
                <div>
                  <h3 class="font-semibold text-gray-900">
                    Healthcare & Services
                  </h3>
                  <p class="text-gray-600">
                    Appointment booking, patient forms, and feedback
                  </p>
                </div>
              </li>
            </ul>
          </div>

          <div
            class="bg-gradient-to-br from-white to-blue-50/30 backdrop-blur-xl border border-white/40 rounded-2xl p-8 shadow-2xl">
            <div class="absolute top-6 right-6">
              <div
                class="h-16 w-16 bg-gradient-to-br from-blue-500/20 to-purple-500/20 rounded-2xl flex items-center justify-center">
                <Globe2 class="w-8 h-8 text-blue-500" />
              </div>
            </div>
            <div class="space-y-6">
              <div
                class="bg-gradient-to-r from-white to-green-50/50 backdrop-blur-sm rounded-xl p-6 border border-white/50 shadow-lg">
                <div class="flex items-center justify-between mb-4">
                  <span
                    class="text-sm font-semibold text-gray-600 uppercase tracking-wide"
                    >Customer Satisfaction</span>
                  <span
                    class="text-3xl font-bold bg-gradient-to-r from-green-600 to-emerald-600 bg-clip-text text-transparent"
                    >+34%</span>
                </div>
                <div
                  class="w-full bg-gradient-to-r from-gray-200 to-gray-100 rounded-full h-3">
                  <div
                    class="bg-gradient-to-r from-green-500 to-emerald-500 h-3 rounded-full shadow-sm"
                    style="width: 84%">
                  </div>
                </div>
              </div>

              <div
                class="bg-gradient-to-r from-white to-blue-50/50 backdrop-blur-sm rounded-xl p-6 border border-white/50 shadow-lg">
                <div class="flex items-center justify-between mb-4">
                  <span
                    class="text-sm font-semibold text-gray-600 uppercase tracking-wide"
                    >Response Rate</span>
                  <span
                    class="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent"
                    >+52%</span>
                </div>
                <div
                  class="w-full bg-gradient-to-r from-gray-200 to-gray-100 rounded-full h-3">
                  <div
                    class="bg-gradient-to-r from-blue-500 to-purple-500 h-3 rounded-full shadow-sm"
                    style="width: 92%">
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Pricing Preview -->
    <section
      id="pricing"
      class="py-20 relative"
      aria-labelledby="pricing-heading">
      <div
        class="absolute inset-0 bg-gradient-to-r from-gray-50/50 to-purple-50/30"
        aria-hidden="true">
      </div>
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <header class="text-center mb-16">
          <h2
            id="pricing-heading"
            class="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">
            Simple, transparent pricing
          </h2>
          <p class="text-xl text-gray-600">
            Choose the plan that fits your business. Upgrade or downgrade
            anytime.
          </p>
        </header>

        <div class="max-w-5xl mx-auto">
          {#if data.plans && data.plans.length > 0}
            <PlanSelector
              plans={data.plans}
              onSelectPlan={handleSelectPlan}
              actionLabel={() => 'Get started'}
              showCurrentBadge={false}
              columns={3} />
          {:else}
            <!-- Fallback static plans if API fails -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div
                class="bg-gradient-to-br from-white to-blue-50/30 backdrop-blur-xl border border-white/40 rounded-2xl p-8 shadow-xl">
                <h3 class="text-xl font-bold text-gray-900 mb-2">Starter</h3>
                <div class="mb-6">
                  <span
                    class="text-4xl font-bold bg-gradient-to-r from-blue-600 to-blue-700 bg-clip-text text-transparent"
                    >$29</span>
                  <span class="text-gray-600">/month</span>
                </div>
                <ul class="space-y-3 mb-8">
                  <li class="flex items-center">
                    <div
                      class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                      <Check class="w-3 h-3 text-white" />
                    </div>
                    <span class="text-gray-600">Up to 2 team members</span>
                  </li>
                  <li class="flex items-center">
                    <div
                      class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                      <Check class="w-3 h-3 text-white" />
                    </div>
                    <span class="text-gray-600">1 organization</span>
                  </li>
                  <li class="flex items-center">
                    <div
                      class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                      <Check class="w-3 h-3 text-white" />
                    </div>
                    <span class="text-gray-600">Unlimited QR codes</span>
                  </li>
                  <li class="flex items-center">
                    <div
                      class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                      <Check class="w-3 h-3 text-white" />
                    </div>
                    <span class="text-gray-600">Analytics & insights</span>
                  </li>
                </ul>
                <Button variant="outline" href="/register" class="w-full">
                  Get started
                </Button>
              </div>

              <div class="relative">
                <div
                  class="absolute -top-4 left-0 right-0 mx-auto w-fit px-4 py-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white text-sm font-semibold rounded-full shadow-lg">
                  Most Popular
                </div>
                <div
                  class="bg-gradient-to-br from-white to-purple-50/30 backdrop-blur-xl border-2 border-blue-500/30 rounded-2xl p-8 shadow-2xl">
                  <h3 class="text-xl font-bold text-gray-900 mb-2">
                    Professional
                  </h3>
                  <div class="mb-6">
                    <span
                      class="text-4xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent"
                      >$79</span>
                    <span class="text-gray-600">/month</span>
                  </div>
                  <ul class="space-y-3 mb-8">
                    <li class="flex items-center">
                      <div
                        class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                        <Check class="w-3 h-3 text-white" />
                      </div>
                      <span class="text-gray-600">Up to 5 team members</span>
                    </li>
                    <li class="flex items-center">
                      <div
                        class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                        <Check class="w-3 h-3 text-white" />
                      </div>
                      <span class="text-gray-600">3 organizations</span>
                    </li>
                    <li class="flex items-center">
                      <div
                        class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                        <Check class="w-3 h-3 text-white" />
                      </div>
                      <span class="text-gray-600">Unlimited QR codes</span>
                    </li>
                    <li class="flex items-center">
                      <div
                        class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                        <Check class="w-3 h-3 text-white" />
                      </div>
                      <span class="text-gray-600">Analytics & insights</span>
                    </li>
                  </ul>
                  <Button variant="gradient" href="/register" class="w-full">
                    Get started
                  </Button>
                </div>
              </div>

              <div
                class="bg-gradient-to-br from-white to-purple-50/30 backdrop-blur-xl border border-white/40 rounded-2xl p-8 shadow-xl">
                <h3 class="text-xl font-bold text-gray-900 mb-2">Enterprise</h3>
                <div class="mb-6">
                  <span
                    class="text-4xl font-bold bg-gradient-to-r from-purple-600 to-indigo-600 bg-clip-text text-transparent"
                    >$199</span>
                  <span class="text-gray-600">/month</span>
                </div>
                <ul class="space-y-3 mb-8">
                  <li class="flex items-center">
                    <div
                      class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                      <Check class="w-3 h-3 text-white" />
                    </div>
                    <span class="text-gray-600">Unlimited team members</span>
                  </li>
                  <li class="flex items-center">
                    <div
                      class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                      <Check class="w-3 h-3 text-white" />
                    </div>
                    <span class="text-gray-600">Unlimited organizations</span>
                  </li>
                  <li class="flex items-center">
                    <div
                      class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                      <Check class="w-3 h-3 text-white" />
                    </div>
                    <span class="text-gray-600">Unlimited QR codes</span>
                  </li>
                  <li class="flex items-center">
                    <div
                      class="h-5 w-5 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full flex items-center justify-center mr-3">
                      <Check class="w-3 h-3 text-white" />
                    </div>
                    <span class="text-gray-600">Analytics & insights</span>
                  </li>
                </ul>
                <Button variant="outline" href="/register" class="w-full">
                  Get started
                </Button>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </section>

    <!-- Contact Section -->
    <section id="contact" class="py-32" aria-labelledby="contact-heading">
      <div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h2 id="contact-heading" class="text-4xl font-bold text-gray-900 mb-6">
          Questions? Let's talk.
        </h2>

        <div class="space-y-8">
          <div>
            <p class="text-xl text-gray-600 mb-8">hello@kyooar.com</p>

            <a
              href="mailto:hello@kyooar.com?subject=Kyooar%20Inquiry"
              class="inline-block"
              aria-label="Send email to Kyooar support team">
              <Button variant="primary" size="lg">Send us an email</Button>
            </a>
          </div>

          <p class="text-gray-500 text-sm">
            We'll get back to you within 24 hours
          </p>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer
      class="relative bg-gradient-to-r from-gray-900/95 to-gray-800/95 backdrop-blur-xl text-gray-300 py-12">
      <div
        class="absolute inset-0 bg-gradient-to-r from-blue-900/10 to-purple-900/10">
      </div>
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div
          class="flex flex-col md:flex-row justify-between items-center space-y-6 md:space-y-0">
          <Logo size="md" variant="white" />

          <div class="flex items-center space-x-8">
            <a
              href="/contact"
              class="text-gray-300 hover:text-white transition-colors duration-150"
              >Contact</a>
            <a
              href="/privacy"
              class="text-gray-300 hover:text-white transition-colors duration-150"
              >Privacy</a>
            <a
              href="/terms"
              class="text-gray-300 hover:text-white transition-colors duration-150"
              >Terms</a>
          </div>
        </div>

        <div class="border-t border-gray-700/50 mt-8 pt-8 text-center text-sm">
          <p>&copy; 2024 Kyooar. All rights reserved.</p>
        </div>
      </div>
    </footer>
  </main>
</div>

<style>
  /* Smooth scrolling for the entire page */
  html {
    scroll-behavior: smooth;
  }

  /* Offset for fixed navbar */
  section {
    scroll-margin-top: 100px;
  }

  /* Custom animations */
  @keyframes float {
    0%,
    100% {
      transform: translateY(0px) rotate(0deg);
    }
    50% {
      transform: translateY(-20px) rotate(10deg);
    }
  }

  @keyframes floatReverse {
    0%,
    100% {
      transform: translateY(0px) rotate(0deg);
    }
    50% {
      transform: translateY(15px) rotate(-8deg);
    }
  }

  @keyframes floatUp {
    0% {
      transform: translateY(0px);
      opacity: 0.3;
    }
    50% {
      opacity: 0.8;
    }
    100% {
      transform: translateY(-100px);
      opacity: 0;
    }
  }

  @keyframes wiggle {
    0%,
    100% {
      transform: rotate(0deg) scale(1);
    }
    25% {
      transform: rotate(3deg) scale(1.05);
    }
    75% {
      transform: rotate(-3deg) scale(0.95);
    }
  }

  @keyframes pulseGlow {
    0%,
    100% {
      transform: scale(1);
      opacity: 0.1;
    }
    50% {
      transform: scale(1.1);
      opacity: 0.2;
    }
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateY(30px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  /* Original animations */
  .animate-float {
    animation: float 6s ease-in-out infinite;
  }

  .animate-float-delayed {
    animation: float 6s ease-in-out infinite;
    animation-delay: -3s;
  }

  .animate-slide-up {
    animation: slideUp 0.8s ease-out forwards;
  }

  .animate-slide-up-delayed {
    animation: slideUp 0.8s ease-out 0.2s forwards;
    opacity: 0;
  }

  .animate-fade-in-delayed {
    animation: fadeIn 0.8s ease-out 0.4s forwards;
    opacity: 0;
  }

  /* New animation classes */
  .animate-bounce-slow {
    animation: bounce 3s ease-in-out infinite;
  }

  .animate-spin-slow {
    animation: spin 8s linear infinite;
  }

  .animate-pulse-slow {
    animation: pulse 4s ease-in-out infinite;
  }

  .animate-float-reverse {
    animation: floatReverse 5s ease-in-out infinite;
  }

  .animate-wiggle {
    animation: wiggle 4s ease-in-out infinite;
  }

  .animate-float-up {
    animation: floatUp 8s ease-in-out infinite;
  }

  .animate-float-up-delayed {
    animation: floatUp 8s ease-in-out infinite;
    animation-delay: -4s;
  }

  .animate-float-up-slow {
    animation: floatUp 12s ease-in-out infinite;
  }

  .animate-pulse-glow {
    animation: pulseGlow 6s ease-in-out infinite;
  }

  .animate-pulse-glow-delayed {
    animation: pulseGlow 6s ease-in-out infinite;
    animation-delay: -3s;
  }

  /* QR Code Scanning Animations */
  @keyframes qrScan {
    0%,
    100% {
      transform: scale(1) rotate(0deg);
      opacity: 0.2;
    }
    25% {
      transform: scale(1.05) rotate(2deg);
      opacity: 0.3;
    }
    50% {
      transform: scale(0.95) rotate(-1deg);
      opacity: 0.25;
    }
    75% {
      transform: scale(1.02) rotate(1deg);
      opacity: 0.35;
    }
  }

  @keyframes scanLine {
    0% {
      transform: scaleY(0);
      opacity: 0;
    }
    50% {
      transform: scaleY(1);
      opacity: 1;
    }
    100% {
      transform: scaleY(0);
      opacity: 0;
    }
  }

  @keyframes scanPulse {
    0%,
    100% {
      transform: scale(1);
      opacity: 0.3;
    }
    50% {
      transform: scale(1.1);
      opacity: 0.6;
    }
  }

  @keyframes scanBounce {
    0%,
    100% {
      transform: translateY(0px);
      opacity: 0.2;
    }
    50% {
      transform: translateY(-5px);
      opacity: 0.4;
    }
  }

  @keyframes laserSweep {
    0% {
      transform: translateX(-100px);
      opacity: 0;
    }
    50% {
      opacity: 0.6;
    }
    100% {
      transform: translateX(100px);
      opacity: 0;
    }
  }

  .animate-qr-scan {
    animation: qrScan 4s ease-in-out infinite;
  }

  .animate-qr-scan-delayed {
    animation: qrScan 4s ease-in-out infinite;
    animation-delay: -1s;
  }

  .animate-qr-scan-slow {
    animation: qrScan 6s ease-in-out infinite;
    animation-delay: -2s;
  }

  .animate-qr-scan-reverse {
    animation: qrScan 5s ease-in-out infinite reverse;
    animation-delay: -1.5s;
  }

  .animate-scan-line {
    animation: scanLine 2s ease-in-out infinite;
  }

  .animate-scan-line-delayed {
    animation: scanLine 2s ease-in-out infinite;
    animation-delay: -1s;
  }

  .animate-scan-pulse {
    animation: scanPulse 3s ease-in-out infinite;
  }

  .animate-scan-bounce {
    animation: scanBounce 2.5s ease-in-out infinite;
  }

  .animate-laser-sweep {
    animation: laserSweep 5s ease-in-out infinite;
  }

  .animate-laser-sweep-delayed {
    animation: laserSweep 5s ease-in-out infinite;
    animation-delay: -2.5s;
  }

  /* Grid pattern */
  .grid-pattern {
    background-image:
      linear-gradient(rgba(59, 130, 246, 0.05) 1px, transparent 1px),
      linear-gradient(90deg, rgba(59, 130, 246, 0.05) 1px, transparent 1px);
    background-size: 50px 50px;
  }
</style>
