<script lang="ts">
  import { Card } from '$lib/components/ui';
  import {
    QrCode,
    MessageSquare,
    Star,
    TrendingUp,
    MapPin,
    Clock,
    Users,
    AlertCircle,
    MessageCircle,
    ThumbsUp,
  } from 'lucide-svelte';

  interface QRCodeData {
    id: string;
    label: string;
    location?: string;
    scans_count: number;
    feedback_count: number;
    conversion_rate: number;
    last_scan?: string;
    is_active: boolean;
    average_rating?: number;
    response_quality?: number; // percentage of responses with comments
  }

  interface AnalyticsData {
    qr_performance?: QRCodeData[];
    total_qr_scans?: number;
    total_active_codes?: number;
    average_conversion_rate?: number;
    total_feedback?: number;
    average_rating?: number;
  }

  let {
    analyticsData = null,
    feedbacks = [],
    loading = false,
  }: {
    analyticsData?: AnalyticsData | null;
    feedbacks?: any[];
    loading?: boolean;
  } = $props();

  // Calculate QR code metrics with response quality focus
  const qrMetrics = $derived(() => {
    if (!analyticsData) return null;

    const qrData = analyticsData.qr_performance || [];
    const totalScans =
      analyticsData.total_qr_scans ||
      qrData.reduce((sum, qr) => sum + qr.scans_count, 0);
    const totalFeedback =
      analyticsData.total_feedback ||
      qrData.reduce((sum, qr) => sum + qr.feedback_count, 0);
    const activeCount =
      analyticsData.total_active_codes ||
      qrData.filter(qr => qr.is_active).length;
    const avgConversion =
      analyticsData.average_conversion_rate ||
      (totalScans > 0 ? (totalFeedback / totalScans) * 100 : 0);

    // Calculate response quality metrics from feedbacks if available
    const qrDataEnhanced = qrData.map(qr => {
      const qrFeedbacks = feedbacks.filter(f => f.qr_code_id === qr.id);
      const ratings = qrFeedbacks
        .map(f => f.overall_rating || f.rating || 0)
        .filter(r => r > 0);
      const avgRating =
        ratings.length > 0
          ? ratings.reduce((sum, r) => sum + r, 0) / ratings.length
          : 0;
      const withComments = qrFeedbacks.filter(
        f => f.comment && f.comment.trim().length > 0
      ).length;
      const responseQuality =
        qrFeedbacks.length > 0 ? (withComments / qrFeedbacks.length) * 100 : 0;

      return {
        ...qr,
        average_rating: avgRating,
        response_quality: responseQuality,
        feedback_count: qrFeedbacks.length || qr.feedback_count,
      };
    });

    // Sort by response generation (feedback count)
    const sortedByResponses = [...qrDataEnhanced].sort(
      (a, b) => b.feedback_count - a.feedback_count
    );
    const topResponseGenerators = sortedByResponses.slice(0, 5);

    // Find QR codes with quality issues
    const qualityIssues = qrDataEnhanced.filter(
      qr =>
        (qr.average_rating > 0 && qr.average_rating < 3) ||
        (qr.response_quality < 30 && qr.feedback_count > 5)
    );

    return {
      totalScans,
      totalFeedback,
      activeCount,
      avgConversion,
      topResponseGenerators,
      qualityIssues,
      allQRCodes: sortedByResponses,
    };
  });

  function getConversionColor(rate: number): string {
    if (rate >= 80) return 'text-green-600';
    if (rate >= 60) return 'text-blue-600';
    if (rate >= 40) return 'text-yellow-600';
    return 'text-red-600';
  }

  function getConversionBg(rate: number): string {
    if (rate >= 80) return 'bg-green-50 border-green-200';
    if (rate >= 60) return 'bg-blue-50 border-blue-200';
    if (rate >= 40) return 'bg-yellow-50 border-yellow-200';
    return 'bg-red-50 border-red-200';
  }

  function formatLastScan(dateString?: string): string {
    if (!dateString) return 'Never';

    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

    if (diffHours < 1) return 'Just now';
    if (diffHours < 24)
      return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`;
    if (diffDays < 7) return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`;

    return date.toLocaleDateString();
  }
</script>

<div class="space-y-6">
  <!-- QR Code Response Analytics -->
  <Card variant="elevated" padding={false}>
    <div class="p-6">
      <div class="mb-6">
        <h3
          class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
          QR Code Response Performance
        </h3>
        <p class="text-sm text-gray-600 mt-1">
          How each QR code generates quality feedback responses
        </p>
      </div>

      {#if loading}
        <div class="space-y-4">
          {#each Array(5) as _}
            <div class="animate-pulse">
              <div class="bg-gray-100 rounded-lg p-4">
                <div class="flex items-center justify-between">
                  <div class="h-4 bg-gray-200 rounded w-1/3"></div>
                  <div class="h-4 bg-gray-200 rounded w-20"></div>
                </div>
              </div>
            </div>
          {/each}
        </div>
      {:else if qrMetrics() && qrMetrics().allQRCodes.length > 0}
        <div class="space-y-4">
          {#each qrMetrics().allQRCodes as qr}
            {@const hasRating = qr.average_rating > 0}
            {@const ratingColor = hasRating
              ? qr.average_rating >= 4
                ? 'text-green-600'
                : qr.average_rating >= 3
                  ? 'text-yellow-600'
                  : 'text-red-600'
              : 'text-gray-400'}
            <Card
              variant={qr.feedback_count > 0 ? 'minimal' : 'default'}
              padding={false}
              hover
              class="p-4 {qr.feedback_count === 0 ? 'bg-gray-50' : ''}">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <div class="flex items-start gap-3">
                    <div
                      class="h-10 w-10 bg-gradient-to-br from-indigo-500 to-purple-500 rounded-lg flex items-center justify-center shadow-sm">
                      <QrCode class="h-5 w-5 text-white" />
                    </div>
                    <div class="flex-1">
                      <h4 class="font-medium text-gray-900">
                        {qr.label || 'Unnamed QR Code'}
                      </h4>
                      {#if qr.location}
                        <div class="flex items-center gap-1 mt-1">
                          <MapPin class="h-3 w-3 text-gray-400" />
                          <span class="text-sm text-gray-600"
                            >{qr.location}</span>
                        </div>
                      {/if}

                      <!-- Response Metrics -->
                      <div class="flex items-center gap-4 mt-3">
                        <div class="flex items-center gap-1">
                          <MessageSquare class="h-4 w-4 text-blue-500" />
                          <span class="text-sm font-medium text-gray-700"
                            >{qr.feedback_count} responses</span>
                        </div>

                        {#if hasRating}
                          <div class="flex items-center gap-1">
                            <Star class="h-4 w-4 {ratingColor} fill-current" />
                            <span class="text-sm font-medium {ratingColor}"
                              >{qr.average_rating.toFixed(1)}</span>
                          </div>
                        {/if}

                        {#if qr.response_quality != null && qr.feedback_count > 0}
                          <div class="flex items-center gap-1">
                            <MessageCircle class="h-4 w-4 text-purple-500" />
                            <span class="text-sm text-gray-600"
                              >{qr.response_quality.toFixed(0)}% with comments</span>
                          </div>
                        {/if}
                      </div>

                      <!-- Activity Info -->
                      <div
                        class="flex items-center gap-3 mt-2 text-xs text-gray-500">
                        <span>{qr.scans_count} total scans</span>
                        <span>•</span>
                        <span>Last: {formatLastScan(qr.last_scan)}</span>
                        {#if !qr.is_active}
                          <span
                            class="px-2 py-0.5 bg-gray-100 text-gray-600 font-medium rounded"
                            >Inactive</span>
                        {/if}
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Response Quality Indicators -->
                <div class="text-right">
                  <div class="space-y-2">
                    <div>
                      <p class="text-xs text-gray-500 uppercase tracking-wide">
                        Response Rate
                      </p>
                      <p
                        class="text-2xl font-bold {getConversionColor(
                          qr.conversion_rate
                        )}">
                        {qr.conversion_rate.toFixed(0)}%
                      </p>
                    </div>

                    {#if qr.feedback_count > 0}
                      <div class="flex items-center justify-end gap-2">
                        {#if hasRating && qr.average_rating >= 4}
                          <span
                            class="px-2 py-1 bg-green-100 text-green-700 text-xs font-medium rounded-full">
                            High Satisfaction
                          </span>
                        {:else if hasRating && qr.average_rating < 3}
                          <span
                            class="px-2 py-1 bg-red-100 text-red-700 text-xs font-medium rounded-full">
                            Needs Attention
                          </span>
                        {/if}
                      </div>
                    {/if}
                  </div>
                </div>
              </div>

              <!-- Response Quality Alerts -->
              {#if qr.feedback_count > 5}
                {#if qr.response_quality < 30}
                  <div class="mt-3 pt-3 border-t border-amber-200">
                    <div class="flex items-center gap-2 text-sm text-amber-700">
                      <AlertCircle class="h-4 w-4" />
                      <span
                        >Low comment rate - Consider encouraging more detailed
                        feedback</span>
                    </div>
                  </div>
                {/if}

                {#if hasRating && qr.average_rating < 3}
                  <div class="mt-3 pt-3 border-t border-red-200">
                    <div class="flex items-center gap-2 text-sm text-red-700">
                      <AlertCircle class="h-4 w-4" />
                      <span>Low satisfaction ratings from this location</span>
                    </div>
                  </div>
                {/if}
              {:else if qr.feedback_count === 0 && qr.scans_count > 5}
                <div class="mt-3 pt-3 border-t border-gray-200">
                  <div class="flex items-center gap-2 text-sm text-gray-600">
                    <AlertCircle class="h-4 w-4" />
                    <span
                      >No responses yet despite {qr.scans_count} scans - Review QR
                      code messaging</span>
                  </div>
                </div>
              {/if}
            </Card>
          {/each}
        </div>

        <!-- Response Quality Insights -->
        {#if qrMetrics().qualityIssues.length > 0}
          <div class="mt-6 p-4 bg-amber-50 border border-amber-200 rounded-lg">
            <div class="flex items-start gap-3">
              <AlertCircle class="h-5 w-5 text-amber-600 mt-0.5" />
              <div>
                <h4 class="font-medium text-amber-900">
                  Response Quality Insights
                </h4>
                <p class="text-sm text-amber-700 mt-1">
                  {qrMetrics().qualityIssues.length} QR code{qrMetrics()
                    .qualityIssues.length > 1
                    ? 's'
                    : ''}
                  need attention for better response quality:
                </p>
                <ul
                  class="text-sm text-amber-700 mt-2 space-y-1 list-disc list-inside">
                  <li>Low rating locations may need service improvements</li>
                  <li>
                    Low comment rates suggest customers need more encouragement
                    to share details
                  </li>
                  <li>Consider adding incentives for detailed feedback</li>
                  <li>
                    Review the feedback form to ensure it's easy to complete
                  </li>
                </ul>
              </div>
            </div>
          </div>
        {/if}

        <!-- Top Response Generators -->
        {#if qrMetrics().topResponseGenerators.length > 0}
          <div class="mt-6 p-4 bg-green-50 border border-green-200 rounded-lg">
            <div class="flex items-start gap-3">
              <ThumbsUp class="h-5 w-5 text-green-600 mt-0.5" />
              <div>
                <h4 class="font-medium text-green-900">
                  Top Response Generators
                </h4>
                <p class="text-sm text-green-700 mt-1">
                  These QR codes are performing exceptionally well:
                </p>
                <ul class="text-sm text-green-700 mt-2 space-y-1">
                  {#each qrMetrics().topResponseGenerators.slice(0, 3) as qr}
                    <li>
                      <span class="font-medium">{qr.label}</span>
                      {#if qr.location}
                        at {qr.location}{/if}
                      - {qr.feedback_count} responses
                      {#if qr.average_rating > 0}
                        ({qr.average_rating.toFixed(1)} ⭐){/if}
                    </li>
                  {/each}
                </ul>
              </div>
            </div>
          </div>
        {/if}
      {:else}
        <div class="text-center py-8">
          <QrCode class="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <p class="text-gray-500">No QR code data available yet.</p>
          <p class="text-sm text-gray-400 mt-2">
            Create and deploy QR codes to start tracking performance.
          </p>
        </div>
      {/if}
    </div>
  </Card>
</div>
