<script lang="ts">
  import { Button } from '$lib/components/ui';
  import { BarChart3, Download } from 'lucide-svelte';

  interface AnalyticsData {
    total_feedback: number;
    average_rating: number;
    feedback_today: number;
    feedback_this_week: number;
    feedback_this_month: number;
  }

  let {
    analyticsData = null,
    loading = false,
    onexportreport = () => {},
  }: {
    analyticsData?: AnalyticsData | null;
    loading?: boolean;
    onexportreport?: () => void;
  } = $props();
</script>

<div class="mb-8">
  <div
    class="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-6">
    <div class="space-y-3">
      <div class="flex items-center space-x-3">
        <div
          class="h-12 w-12 bg-gradient-to-br from-purple-500 to-pink-600 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25">
          <BarChart3 class="h-6 w-6 text-white" />
        </div>
        <div>
          <h1
            class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
            Analytics & Reports
          </h1>
          <div class="flex items-center space-x-4 mt-1">
            <p class="text-gray-600 font-medium">
              Detailed insights and performance metrics for your organizations
            </p>
            {#if !loading && analyticsData}
              <div class="flex items-center space-x-3 text-sm">
                <div class="flex items-center space-x-1">
                  <div class="w-2 h-2 bg-purple-400 rounded-full"></div>
                  <span class="text-gray-600"
                    >{analyticsData.total_feedback} Total Feedback</span>
                </div>
                <div class="flex items-center space-x-1">
                  <div class="w-2 h-2 bg-pink-400 rounded-full"></div>
                  <span class="text-gray-600"
                    >{analyticsData.average_rating?.toFixed(1) || '0.0'} Avg Rating</span>
                </div>
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>

    <div class="flex items-center space-x-3">
      <!-- Export Report Button -->
      <Button
        variant="gradient"
        size="lg"
        class="group relative overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300"
        onclick={onexportreport}
        disabled={!analyticsData}>
        <div
          class="absolute inset-0 bg-gradient-to-r from-purple-600 to-pink-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
        </div>
        <Download class="h-5 w-5 mr-2 relative z-10 group-hover:scale-110 transition-transform duration-200" />
        <span class="relative z-10">Export Report</span>
      </Button>
    </div>
  </div>
</div>
