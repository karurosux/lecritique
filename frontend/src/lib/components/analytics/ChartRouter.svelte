<script lang="ts">
  import RatingChart from './charts/RatingChart.svelte';
  import ScaleChart from './charts/ScaleChart.svelte';
  import YesNoChart from './charts/YesNoChart.svelte';
  import SentimentChart from './charts/SentimentChart.svelte';
  import ChoiceChart from './charts/ChoiceChart.svelte';
  import CountChart from './charts/CountChart.svelte';
  import ComparisonRouter from './charts/ComparisonRouter.svelte';

  interface Props {
    data: any;
    type: 'timeseries' | 'comparison';
  }

  let { data, type }: Props = $props();

  // Group series by metric type for better visualization
  let groupedSeries = $derived(() => {
    if (!data?.series) return {};

    const groups: Record<string, any[]> = {};
    data.series.forEach((series: any) => {
      const metricType = getMetricGroup(series.metric_type);
      if (!groups[metricType]) {
        groups[metricType] = [];
      }
      groups[metricType].push(series);
    });

    return groups;
  });

  function getMetricGroup(metricType: string): string {
    if (metricType.startsWith('question_')) {
      // Extract question type from metadata or infer from metric name
      return 'question';
    }

    if (metricType.includes('rating')) return 'rating';
    if (metricType.includes('scale')) return 'scale';
    if (metricType.includes('yes_no')) return 'yes_no';
    if (metricType.includes('text')) return 'sentiment';
    if (metricType.includes('choice')) return 'choice';
    if (metricType.includes('survey_responses')) return 'count';

    return 'count'; // default
  }

  function getQuestionType(series: any): string {
    // Try to get question type from metadata
    if (series.metadata) {
      let metadata = series.metadata;
      if (metadata?.question_type) {
        return metadata.question_type;
      }
    }

    // Fallback to inferring from metric_type
    if (series.metric_type.includes('rating')) return 'rating';
    if (series.metric_type.includes('scale')) return 'scale';
    if (series.metric_type.includes('yes_no')) return 'yes_no';
    if (series.metric_type.includes('text')) return 'text';
    if (series.metric_type.includes('choice')) return 'single_choice';

    return 'count';
  }
</script>

<div class="chart-router">
  {#if type === 'comparison'}
    <ComparisonRouter {data} />
  {:else}
    {#each Object.entries(groupedSeries) as [metricGroup, seriesData]}
      <div class="metric-group mb-8">
        {#if metricGroup === 'question'}
          <!-- Handle individual questions based on their type -->
          {#each seriesData as series}
            {@const questionType = getQuestionType(series)}
            <div class="question-chart mb-6">
              {#if questionType === 'rating'}
                <RatingChart data={{ ...data, series: [series] }} {type} />
              {:else if questionType === 'scale'}
                <ScaleChart data={{ ...data, series: [series] }} {type} />
              {:else if questionType === 'yes_no'}
                <YesNoChart data={{ ...data, series: [series] }} {type} />
              {:else if questionType === 'text'}
                <SentimentChart data={{ ...data, series: [series] }} {type} />
              {:else if questionType.includes('choice')}
                <ChoiceChart data={{ ...data, series: [series] }} {type} />
              {:else}
                <CountChart data={{ ...data, series: [series] }} {type} />
              {/if}
            </div>
          {/each}
        {:else if metricGroup === 'rating'}
          <RatingChart data={{ ...data, series: seriesData }} {type} />
        {:else if metricGroup === 'scale'}
          <ScaleChart data={{ ...data, series: seriesData }} {type} />
        {:else if metricGroup === 'yes_no'}
          <YesNoChart data={{ ...data, series: seriesData }} {type} />
        {:else if metricGroup === 'sentiment'}
          <SentimentChart data={{ ...data, series: seriesData }} {type} />
        {:else if metricGroup === 'choice'}
          <ChoiceChart data={{ ...data, series: seriesData }} {type} />
        {:else}
          <CountChart data={{ ...data, series: seriesData }} {type} />
        {/if}
      </div>
    {/each}
  {/if}
</div>
