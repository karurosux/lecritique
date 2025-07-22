<script lang="ts">
  import { Card } from '$lib/components/ui';
  import { onMount } from 'svelte';
  
  interface QuestionMetric {
    question_text: string;
    average_score?: number;
    response_count: number;
  }
  
  interface ProductInsights {
    product_name?: string;
    question_metrics: QuestionMetric[];
  }
  
  let {
    productInsights = null,
    loading = false
  }: {
    productInsights?: ProductInsights | null;
    loading?: boolean;
  } = $props();

  let canvas: HTMLCanvasElement;
  let ctx: CanvasRenderingContext2D | null;
  let animationProgress = $state(0);

  const chartData = $derived(() => {
    if (!productInsights?.question_metrics) return null;
    
    const numericQuestions = productInsights.question_metrics.filter(
      q => q.average_score !== undefined && q.average_score !== null
    );
    
    if (numericQuestions.length < 3) return null;
    
    return {
      labels: numericQuestions.map(q => {
        const text = q.question_text;
        return text.length > 20 ? text.substring(0, 17) + '...' : text;
      }),
      scores: numericQuestions.map(q => q.average_score || 0),
      maxScore: 5
    };
  });

  function drawRadarChart() {
    if (!canvas || !ctx || !chartData) return;
    
    const centerX = canvas.width / 2;
    const centerY = canvas.height / 2;
    const radius = Math.min(centerX, centerY) - 60;
    const angleStep = (Math.PI * 2) / chartData.labels.length;
    
    // Clear canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    // Draw grid circles
    for (let i = 1; i <= 5; i++) {
      ctx.beginPath();
      ctx.strokeStyle = '#e5e7eb';
      ctx.lineWidth = 1;
      
      for (let j = 0; j <= chartData.labels.length; j++) {
        const angle = j * angleStep - Math.PI / 2;
        const x = centerX + Math.cos(angle) * (radius * i / 5);
        const y = centerY + Math.sin(angle) * (radius * i / 5);
        
        if (j === 0) {
          ctx.moveTo(x, y);
        } else {
          ctx.lineTo(x, y);
        }
      }
      ctx.closePath();
      ctx.stroke();
    }
    
    // Draw axes
    chartData.labels.forEach((_, index) => {
      const angle = index * angleStep - Math.PI / 2;
      ctx.beginPath();
      ctx.strokeStyle = '#d1d5db';
      ctx.lineWidth = 1;
      ctx.moveTo(centerX, centerY);
      ctx.lineTo(
        centerX + Math.cos(angle) * radius,
        centerY + Math.sin(angle) * radius
      );
      ctx.stroke();
    });
    
    // Draw data polygon with animation
    ctx.beginPath();
    ctx.fillStyle = 'rgba(59, 130, 246, 0.2)';
    ctx.strokeStyle = '#3b82f6';
    ctx.lineWidth = 2;
    
    chartData.scores.forEach((score, index) => {
      const angle = index * angleStep - Math.PI / 2;
      const distance = (score / chartData.maxScore) * radius * animationProgress;
      const x = centerX + Math.cos(angle) * distance;
      const y = centerY + Math.sin(angle) * distance;
      
      if (index === 0) {
        ctx.moveTo(x, y);
      } else {
        ctx.lineTo(x, y);
      }
    });
    
    ctx.closePath();
    ctx.fill();
    ctx.stroke();
    
    // Draw data points
    chartData.scores.forEach((score, index) => {
      const angle = index * angleStep - Math.PI / 2;
      const distance = (score / chartData.maxScore) * radius * animationProgress;
      const x = centerX + Math.cos(angle) * distance;
      const y = centerY + Math.sin(angle) * distance;
      
      ctx.beginPath();
      ctx.fillStyle = '#3b82f6';
      ctx.arc(x, y, 4, 0, Math.PI * 2);
      ctx.fill();
    });
    
    // Draw labels
    ctx.font = '12px system-ui, -apple-system, sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    
    chartData.labels.forEach((label, index) => {
      const angle = index * angleStep - Math.PI / 2;
      const labelDistance = radius + 30;
      const x = centerX + Math.cos(angle) * labelDistance;
      const y = centerY + Math.sin(angle) * labelDistance;
      
      ctx.fillStyle = '#374151';
      ctx.fillText(label, x, y);
      
      // Draw score value
      const score = chartData.scores[index];
      const scoreDistance = (score / chartData.maxScore) * radius * animationProgress + 15;
      const scoreX = centerX + Math.cos(angle) * scoreDistance;
      const scoreY = centerY + Math.sin(angle) * scoreDistance;
      
      ctx.fillStyle = '#1f2937';
      ctx.font = 'bold 11px system-ui, -apple-system, sans-serif';
      ctx.fillText(score.toFixed(1), scoreX, scoreY);
    });
  }

  onMount(() => {
    ctx = canvas?.getContext('2d');
    
    // Animate chart
    const animate = () => {
      if (animationProgress < 1) {
        animationProgress = Math.min(animationProgress + 0.02, 1);
        drawRadarChart();
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
          drawRadarChart();
          requestAnimationFrame(animate);
        }
      };
      animate();
    }
  });
</script>

<Card variant="default" class="product-radar-chart">
  <div class="mb-4">
    <h3 class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
      Question Score Radar
    </h3>
    <p class="text-sm text-gray-600 mt-1">
      Visual comparison of scores across all questions
    </p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center h-[400px]">
      <div class="animate-pulse">
        <div class="w-64 h-64 bg-gray-200 rounded-full"></div>
      </div>
    </div>
  {:else if chartData}
    <div class="relative">
      <canvas 
        bind:this={canvas}
        width="500"
        height="400"
        class="w-full max-w-[500px] mx-auto"
      ></canvas>
    </div>
  {:else}
    <div class="text-center py-12">
      <div class="text-gray-500 text-sm">
        Not enough numeric questions to display radar chart (minimum 3 required)
      </div>
    </div>
  {/if}
</Card>
