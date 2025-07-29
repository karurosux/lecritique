<script lang="ts">
  import { Card } from '$lib/components/ui';
  import { onMount } from 'svelte';
  import {
    CalendarIcon,
    TrendingUpIcon,
    TrendingDownIcon,
  } from 'lucide-svelte';

  interface TrendData {
    date: string;
    average_rating: number;
    feedback_count: number;
  }

  interface ProductInsights {
    product_name?: string;
    rating_trend?: TrendData[];
    weekly_change?: number;
  }

  let {
    productInsights = null,
    loading = false,
    timeframe = '7d',
  }: {
    productInsights?: ProductInsights | null;
    loading?: boolean;
    timeframe?: string;
  } = $props();

  let canvas: HTMLCanvasElement;
  let ctx: CanvasRenderingContext2D | null;
  let hoveredPoint: number | null = null;
  let animationProgress = $state(0);

  const chartData = $derived(() => {
    if (
      !productInsights?.rating_trend ||
      productInsights.rating_trend.length < 2
    )
      return null;

    const data = productInsights.rating_trend.slice(-30); // Last 30 data points
    return {
      labels: data.map(d => {
        const date = new Date(d.date);
        return date.toLocaleDateString('en-US', {
          month: 'short',
          day: 'numeric',
        });
      }),
      ratings: data.map(d => d.average_rating || 0),
      counts: data.map(d => d.feedback_count || 0),
      minRating: Math.min(...data.map(d => d.average_rating || 0)),
      maxRating: Math.max(...data.map(d => d.average_rating || 0)),
    };
  });

  function drawTrendChart() {
    if (!canvas || !ctx || !chartData) return;

    const padding = { top: 20, right: 20, bottom: 50, left: 50 };
    const chartWidth = canvas.width - padding.left - padding.right;
    const chartHeight = canvas.height - padding.top - padding.bottom;

    // Clear canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Draw grid lines
    ctx.strokeStyle = '#e5e7eb';
    ctx.lineWidth = 1;

    // Y-axis grid lines
    for (let i = 0; i <= 5; i++) {
      const y = padding.top + (chartHeight / 5) * (5 - i);
      ctx.beginPath();
      ctx.moveTo(padding.left, y);
      ctx.lineTo(canvas.width - padding.right, y);
      ctx.stroke();

      // Y-axis labels
      ctx.fillStyle = '#6b7280';
      ctx.font = '11px system-ui, -apple-system, sans-serif';
      ctx.textAlign = 'right';
      ctx.textBaseline = 'middle';
      ctx.fillText(i.toString(), padding.left - 10, y);
    }

    // Draw data line
    const xStep = chartWidth / (chartData.labels.length - 1);

    // Gradient fill
    const gradient = ctx.createLinearGradient(
      0,
      padding.top,
      0,
      canvas.height - padding.bottom
    );
    gradient.addColorStop(0, 'rgba(59, 130, 246, 0.3)');
    gradient.addColorStop(1, 'rgba(59, 130, 246, 0.0)');

    // Draw filled area
    ctx.beginPath();
    ctx.fillStyle = gradient;

    chartData.ratings.forEach((rating, index) => {
      const x = padding.left + index * xStep;
      const progress =
        Math.min(animationProgress * chartData.ratings.length, index + 1) /
        (index + 1);
      const y = padding.top + (1 - (rating / 5) * progress) * chartHeight;

      if (index === 0) {
        ctx.moveTo(x, y);
      } else {
        ctx.lineTo(x, y);
      }
    });

    // Complete the fill
    ctx.lineTo(
      padding.left + (chartData.ratings.length - 1) * xStep,
      canvas.height - padding.bottom
    );
    ctx.lineTo(padding.left, canvas.height - padding.bottom);
    ctx.closePath();
    ctx.fill();

    // Draw line
    ctx.beginPath();
    ctx.strokeStyle = '#3b82f6';
    ctx.lineWidth = 2;

    chartData.ratings.forEach((rating, index) => {
      const x = padding.left + index * xStep;
      const progress =
        Math.min(animationProgress * chartData.ratings.length, index + 1) /
        (index + 1);
      const y = padding.top + (1 - (rating / 5) * progress) * chartHeight;

      if (index === 0) {
        ctx.moveTo(x, y);
      } else {
        ctx.lineTo(x, y);
      }
    });

    ctx.stroke();

    // Draw data points
    chartData.ratings.forEach((rating, index) => {
      const x = padding.left + index * xStep;
      const progress =
        Math.min(animationProgress * chartData.ratings.length, index + 1) /
        (index + 1);
      const y = padding.top + (1 - (rating / 5) * progress) * chartHeight;

      ctx.beginPath();
      ctx.fillStyle = hoveredPoint === index ? '#2563eb' : '#3b82f6';
      ctx.arc(x, y, hoveredPoint === index ? 6 : 4, 0, Math.PI * 2);
      ctx.fill();

      // Draw hover tooltip
      if (hoveredPoint === index) {
        const tooltipWidth = 120;
        const tooltipHeight = 50;
        const tooltipX = Math.min(
          x - tooltipWidth / 2,
          canvas.width - tooltipWidth - 10
        );
        const tooltipY = y - tooltipHeight - 10;

        // Tooltip background
        ctx.fillStyle = 'rgba(31, 41, 55, 0.9)';
        ctx.beginPath();
        ctx.moveTo(tooltipX + 4, tooltipY);
        ctx.lineTo(tooltipX + tooltipWidth - 4, tooltipY);
        ctx.quadraticCurveTo(
          tooltipX + tooltipWidth,
          tooltipY,
          tooltipX + tooltipWidth,
          tooltipY + 4
        );
        ctx.lineTo(tooltipX + tooltipWidth, tooltipY + tooltipHeight - 4);
        ctx.quadraticCurveTo(
          tooltipX + tooltipWidth,
          tooltipY + tooltipHeight,
          tooltipX + tooltipWidth - 4,
          tooltipY + tooltipHeight
        );
        ctx.lineTo(tooltipX + 4, tooltipY + tooltipHeight);
        ctx.quadraticCurveTo(
          tooltipX,
          tooltipY + tooltipHeight,
          tooltipX,
          tooltipY + tooltipHeight - 4
        );
        ctx.lineTo(tooltipX, tooltipY + 4);
        ctx.quadraticCurveTo(tooltipX, tooltipY, tooltipX + 4, tooltipY);
        ctx.closePath();
        ctx.fill();

        // Tooltip text
        ctx.fillStyle = 'white';
        ctx.font = 'bold 12px system-ui, -apple-system, sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'top';
        ctx.fillText(
          chartData.labels[index],
          tooltipX + tooltipWidth / 2,
          tooltipY + 8
        );

        ctx.font = '11px system-ui, -apple-system, sans-serif';
        ctx.fillText(
          `Rating: ${rating.toFixed(1)}`,
          tooltipX + tooltipWidth / 2,
          tooltipY + 24
        );
        ctx.fillText(
          `Reviews: ${chartData.counts[index]}`,
          tooltipX + tooltipWidth / 2,
          tooltipY + 36
        );
      }
    });

    // Draw x-axis labels
    ctx.fillStyle = '#6b7280';
    ctx.font = '10px system-ui, -apple-system, sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'top';

    const labelStep = Math.ceil(chartData.labels.length / 8);
    chartData.labels.forEach((label, index) => {
      if (index % labelStep === 0 || index === chartData.labels.length - 1) {
        const x = padding.left + index * xStep;
        ctx.fillText(label, x, canvas.height - padding.bottom + 10);
      }
    });
  }

  function handleMouseMove(event: MouseEvent) {
    if (!canvas || !chartData) return;

    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const padding = { left: 50, right: 20 };
    const chartWidth = canvas.width - padding.left - padding.right;
    const xStep = chartWidth / (chartData.labels.length - 1);

    const index = Math.round((x - padding.left) / xStep);
    if (
      index >= 0 &&
      index < chartData.labels.length &&
      x >= padding.left &&
      x <= canvas.width - padding.right
    ) {
      hoveredPoint = index;
    } else {
      hoveredPoint = null;
    }

    drawTrendChart();
  }

  function handleMouseLeave() {
    hoveredPoint = null;
    drawTrendChart();
  }

  onMount(() => {
    ctx = canvas?.getContext('2d');

    // Animate chart
    const animate = () => {
      if (animationProgress < 1) {
        animationProgress = Math.min(animationProgress + 0.02, 1);
        drawTrendChart();
        requestAnimationFrame(animate);
      }
    };

    if (chartData) {
      animate();
    }
  });

  $effect(() => {
    if (ctx && chartData) {
      animationProgress = 0;
      const animate = () => {
        if (animationProgress < 1) {
          animationProgress = Math.min(animationProgress + 0.02, 1);
          drawTrendChart();
          requestAnimationFrame(animate);
        }
      };
      animate();
    }
  });
</script>

<Card variant="default" class="product-trend-chart">
  <div class="mb-4 flex items-center justify-between">
    <div>
      <h3
        class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
        Rating Trend
      </h3>
      <p class="text-sm text-gray-600 mt-1">
        {productInsights?.product_name || 'Product'} performance over time
      </p>
    </div>

    {#if productInsights?.weekly_change !== undefined}
      <div class="flex items-center gap-2">
        {#if productInsights.weekly_change > 0}
          <TrendingUpIcon class="w-4 h-4 text-green-600" />
          <span class="text-sm font-medium text-green-600">
            +{productInsights.weekly_change.toFixed(1)}% this week
          </span>
        {:else if productInsights.weekly_change < 0}
          <TrendingDownIcon class="w-4 h-4 text-red-600" />
          <span class="text-sm font-medium text-red-600">
            {productInsights.weekly_change.toFixed(1)}% this week
          </span>
        {:else}
          <span class="text-sm font-medium text-gray-600">
            No change this week
          </span>
        {/if}
      </div>
    {/if}
  </div>

  {#if loading}
    <div class="flex items-center justify-center h-[300px]">
      <div class="animate-pulse w-full h-full bg-gray-100 rounded"></div>
    </div>
  {:else if chartData}
    <div class="relative">
      <canvas
        bind:this={canvas}
        width="600"
        height="300"
        class="w-full cursor-crosshair"
        onmousemove={handleMouseMove}
        onmouseleave={handleMouseLeave}></canvas>

      <div
        class="mt-4 flex items-center justify-center gap-4 text-xs text-gray-500">
        <div class="flex items-center gap-1">
          <CalendarIcon class="w-3 h-3" />
          <span>Last {chartData.labels.length} days</span>
        </div>
        <div class="flex items-center gap-1">
          <div class="w-3 h-3 bg-blue-500 rounded-full"></div>
          <span>Average Rating</span>
        </div>
      </div>
    </div>
  {:else}
    <div class="text-center py-12">
      <div class="text-gray-500 text-sm">
        Not enough trend data available. Need at least 2 days of feedback.
      </div>
    </div>
  {/if}
</Card>
