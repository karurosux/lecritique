<script lang="ts">
  import ComparisonChart from '../ComparisonChart.svelte';
  
  interface Props {
    data: any;
  }
  
  let { data }: Props = $props();
  
  // Group comparisons by metric type for better visualization
  let groupedComparisons = $derived(() => {
    if (!data?.comparisons) return {};
    
    const groups: Record<string, any[]> = {};
    data.comparisons.forEach((comparison: any) => {
      const metricType = getMetricGroup(comparison.metric_type);
      if (!groups[metricType]) {
        groups[metricType] = [];
      }
      groups[metricType].push(comparison);
    });
    
    return groups;
  });
  
  function getMetricGroup(metricType: string): string {
    if (metricType.startsWith('question_')) {
      return 'question';
    }
    
    if (metricType.includes('rating')) return 'rating';
    if (metricType.includes('scale')) return 'scale'; 
    if (metricType.includes('yes_no')) return 'yes_no';
    if (metricType.includes('text')) return 'sentiment';
    if (metricType.includes('choice')) return 'choice';
    if (metricType.includes('survey_responses')) return 'count';
    
    return 'count';
  }
  
  function getQuestionType(comparison: any): string {
    if (comparison.metadata) {
      let metadata = comparison.metadata;
      if (typeof metadata === 'string') {
        try {
          metadata = JSON.parse(metadata);
        } catch (e) {
          // ignore parsing errors
        }
      }
      if (metadata?.question_type) {
        return metadata.question_type;
      }
    }
    
    if (comparison.metric_type.includes('rating')) return 'rating';
    if (comparison.metric_type.includes('scale')) return 'scale';
    if (comparison.metric_type.includes('yes_no')) return 'yes_no';
    if (comparison.metric_type.includes('text')) return 'text';
    if (comparison.metric_type.includes('choice')) return 'single_choice';
    
    return 'count';
  }
</script>

<div class="comparison-router">
  <!-- For now, use the existing ComparisonChart but with better grouping coming soon -->
  <ComparisonChart {data} />
</div>