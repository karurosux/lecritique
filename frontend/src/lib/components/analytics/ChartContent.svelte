<script lang="ts">
  import {
    Star,
    BarChart3,
    ToggleLeft,
    ToggleRight,
    MessageSquare,
    ThumbsUp,
    ThumbsDown,
    Meh,
  } from "lucide-svelte";

  interface BackendChartData {
    question_id: string;
    question_text: string;
    question_type: string;
    chart_type: string;
    data: {
      scale?: number;
      distribution?: Record<string, number>;
      average?: number;
      total?: number;
      percentages?: Record<string, number>;
      options?: Record<string, number>;
      is_multi_choice?: boolean;
      combinations?: Array<{
        options: string[];
        count: number;
        percentage: number;
      }>;
      positive?: number;
      neutral?: number;
      negative?: number;
      samples?: string[];
      keywords?: string[];
    };
  }

  let {
    chart,
    chartViewModes = $bindable(new Map())
  }: {
    chart: BackendChartData;
    chartViewModes?: Map<string, "individual" | "combinations">;
  } = $props();

  function getQuestionTypeLabel(type: string): string {
    switch (type) {
      case "rating":
        return "Rating";
      case "scale":
        return "Scale";
      case "single_choice":
        return "Single Choice";
      case "multi_choice":
        return "Multiple Choice";
      case "yes_no":
        return "Yes/No";
      case "text":
        return "Text Response";
      default:
        return "Response";
    }
  }

  function getRatingStars(rating: number, scale: number = 5): string {
    const normalizedRating = scale === 10 ? rating / 2 : rating;
    const fullStars = Math.floor(normalizedRating);
    const hasHalfStar = normalizedRating % 1 >= 0.5;

    let stars = "";
    for (let i = 0; i < fullStars; i++) {
      stars += "★";
    }
    if (hasHalfStar) {
      stars += "☆";
    }
    while (stars.length < 5) {
      stars += "☆";
    }
    return stars;
  }

  function getColorForRating(rating: number, scale: number = 5): string {
    const normalizedRating = scale === 10 ? rating / 2 : rating;
    if (normalizedRating >= 4) return "text-emerald-600";
    if (normalizedRating >= 3) return "text-yellow-600";
    return "text-red-600";
  }

  function getColorForChoice(index: number): string {
    const colors = [
      "bg-gradient-to-r from-blue-500 to-blue-600",
      "bg-gradient-to-r from-emerald-500 to-emerald-600",
      "bg-gradient-to-r from-purple-500 to-purple-600",
      "bg-gradient-to-r from-orange-500 to-orange-600",
      "bg-gradient-to-r from-pink-500 to-pink-600",
      "bg-gradient-to-r from-indigo-500 to-indigo-600",
      "bg-gradient-to-r from-teal-500 to-teal-600",
      "bg-gradient-to-r from-amber-500 to-amber-600",
    ];
    return colors[index % colors.length];
  }

  function toggleViewMode(questionId: string) {
    const current = chartViewModes.get(questionId) || "individual";
    chartViewModes.set(
      questionId,
      current === "individual" ? "combinations" : "individual",
    );
    chartViewModes = new Map(chartViewModes);
  }
</script>

<!-- Chart Header -->
<div class="flex items-start justify-between mb-6">
  <div class="flex-1">
    <div class="flex items-center gap-3 mb-3">
      <div
        class="h-8 w-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center"
      >
        {#if chart.chart_type === "rating"}
          <Star class="h-4 w-4 text-white" />
        {:else if chart.chart_type === "scale"}
          <BarChart3 class="h-4 w-4 text-white" />
        {:else if chart.chart_type === "text_sentiment"}
          <MessageSquare class="h-4 w-4 text-white" />
        {:else}
          <BarChart3 class="h-4 w-4 text-white" />
        {/if}
      </div>
      <div class="flex-1">
        <h3
          class="text-lg font-semibold text-gray-900 leading-tight mb-1"
        >
          {chart.question_text}
        </h3>
        <div class="flex items-center gap-2">
          <span
            class="px-2 py-1 bg-gray-100 text-gray-700 text-xs font-medium rounded-md"
          >
            {getQuestionTypeLabel(chart.question_type)}
          </span>
          <span class="text-sm text-gray-500">
            {(chart.data.total || 0).toLocaleString()} responses
          </span>
        </div>
      </div>
    </div>
  </div>

  <!-- Multi-choice toggle -->
  {#if chart.data.is_multi_choice && chart.data.combinations}
    <button
      onclick={() => toggleViewMode(chart.question_id)}
      class="flex items-center gap-2 px-3 py-1.5 text-sm font-medium text-gray-600 hover:text-gray-900 bg-gray-50 hover:bg-gray-100 rounded-lg border border-gray-200 transition-all duration-200"
    >
      <span class="text-xs">
        {chartViewModes.get(chart.question_id) === "combinations"
          ? "Individual"
          : "Combinations"}
      </span>
      {#if chartViewModes.get(chart.question_id) === "combinations"}
        <ToggleRight class="h-3 w-3 text-emerald-600" />
      {:else}
        <ToggleLeft class="h-3 w-3 text-gray-400" />
      {/if}
    </button>
  {/if}
</div>

<!-- Chart Content -->
{#if chart.chart_type === "rating"}
  <!-- Rating Distribution Chart -->
  <div class="space-y-4">
    {#if chart.data.average}
      <div
        class="text-center mb-6 p-5 bg-gradient-to-br from-amber-50 via-orange-50 to-red-50 rounded-2xl border border-amber-200 relative overflow-hidden"
      >
        <!-- Food pattern background -->
        <div class="absolute inset-0 opacity-10">
          <div class="absolute top-2 right-4 text-2xl">🍽️</div>
          <div class="absolute bottom-2 left-4 text-xl">✨</div>
          <div class="absolute top-4 left-8 text-lg">👨‍🍳</div>
        </div>
        <div class="relative z-10">
          <div
            class="text-3xl font-black text-orange-700 mb-2 tracking-tight"
          >
            {chart.data.average.toFixed(1)}
            <span class="text-lg text-orange-600 ml-1"
              >/{chart.data.scale || 5}</span
            >
          </div>
          <div class="text-2xl mb-3 tracking-wider">
            {getRatingStars(chart.data.average, chart.data.scale)}
          </div>
          <div
            class="inline-flex items-center gap-2 bg-white/80 backdrop-blur-sm px-3 py-1 rounded-full border border-orange-200"
          >
            <Star class="h-3 w-3 text-orange-600" />
            <span class="text-xs font-semibold text-orange-700"
              >Average Rating</span
            >
          </div>
        </div>
      </div>
    {/if}

    <!-- Rating bars -->
    {#if chart.data.distribution}
      <div class="space-y-3">
        {#each Object.entries(chart.data.distribution).sort(([a], [b]) => parseInt(b) - parseInt(a)) as [rating, count], index}
          {@const percentage = chart.data.percentages?.[rating] || 0}
          <div
            class="group p-3 hover:bg-orange-50 rounded-xl border border-orange-100 hover:border-orange-200 transition-all duration-200"
          >
            <div class="flex items-center gap-4">
              <div
                class="flex items-center gap-3 w-20 flex-shrink-0"
              >
                <div
                  class="w-8 h-8 bg-gradient-to-br from-orange-100 to-orange-200 rounded-lg flex items-center justify-center font-bold text-orange-800 text-sm group-hover:scale-105 transition-transform"
                >
                  {rating}
                </div>
                <Star
                  class="h-4 w-4 {getColorForRating(
                    parseInt(rating),
                    chart.data.scale,
                  )} group-hover:scale-105 transition-transform"
                />
              </div>
              <div class="flex-1 relative">
                <div
                  class="bg-gradient-to-r from-gray-200 to-gray-100 rounded-xl h-6 overflow-hidden shadow-inner"
                >
                  <div
                    class="h-full bg-gradient-to-r from-orange-400 via-orange-500 to-orange-600 transition-all duration-700 ease-out shadow-sm rounded-xl relative"
                    style="width: {percentage}%"
                  >
                    <!-- Subtle shine effect -->
                    <div
                      class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent rounded-xl"
                    ></div>
                  </div>
                </div>
                <!-- Percentage overlay for larger bars -->
                {#if percentage > 20}
                  <div
                    class="absolute inset-0 flex items-center justify-center text-xs font-bold text-white pointer-events-none"
                    style="text-shadow: 0 1px 2px rgba(0,0,0,0.5), 0 0 4px rgba(0,0,0,0.3);"
                  >
                    {percentage.toFixed(0)}%
                  </div>
                {/if}
              </div>
              <div
                class="text-sm font-bold text-gray-900 w-20 text-right bg-white border border-orange-200 px-3 py-2 rounded-lg group-hover:bg-orange-50 transition-colors"
              >
                {count.toLocaleString()}
                <span
                  class="block text-xs text-orange-600 font-normal"
                >
                  {percentage.toFixed(1)}%
                </span>
              </div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
{:else if chart.chart_type === "scale"}
  <!-- Scale Distribution Chart -->
  <div class="space-y-4">
    {#if chart.data.average}
      <div
        class="text-center mb-6 p-5 bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 rounded-2xl border border-blue-200 relative overflow-hidden"
      >
        <!-- Scale pattern background -->
        <div class="absolute inset-0 opacity-10">
          <div class="absolute top-2 right-4 text-2xl">📊</div>
          <div class="absolute bottom-2 left-4 text-xl">⚖️</div>
          <div class="absolute top-4 left-8 text-lg">📈</div>
        </div>
        <div class="relative z-10">
          <div
            class="text-3xl font-black text-blue-700 mb-2 tracking-tight"
          >
            {chart.data.average.toFixed(1)}
            <span class="text-lg text-blue-600 ml-1"
              >/{chart.data.scale || 5}</span
            >
          </div>
          <div
            class="inline-flex items-center gap-2 bg-white/80 backdrop-blur-sm px-3 py-1 rounded-full border border-blue-200"
          >
            <BarChart3 class="h-3 w-3 text-blue-600" />
            <span class="text-xs font-semibold text-blue-700"
              >Average Score</span
            >
          </div>
        </div>
      </div>
    {/if}

    <!-- Scale bars -->
    {#if chart.data.distribution}
      <div class="space-y-3">
        {#each Object.entries(chart.data.distribution).sort(([a], [b]) => parseInt(a) - parseInt(b)) as [value, count], index}
          {@const percentage = chart.data.percentages?.[value] || 0}
          <div
            class="group p-3 hover:bg-blue-50 rounded-xl border border-blue-100 hover:border-blue-200 transition-all duration-200"
          >
            <div class="flex items-center gap-4">
              <div
                class="flex items-center gap-3 w-20 flex-shrink-0"
              >
                <div
                  class="w-8 h-8 bg-gradient-to-br from-blue-100 to-blue-200 rounded-lg flex items-center justify-center font-bold text-blue-800 text-sm group-hover:scale-105 transition-transform"
                >
                  {value}
                </div>
                <BarChart3
                  class="h-4 w-4 text-blue-600 group-hover:scale-105 transition-transform"
                />
              </div>
              <div class="flex-1 relative">
                <div
                  class="bg-gradient-to-r from-gray-200 to-gray-100 rounded-xl h-6 overflow-hidden shadow-inner"
                >
                  <div
                    class="h-full bg-gradient-to-r from-blue-400 via-blue-500 to-blue-600 transition-all duration-700 ease-out shadow-sm rounded-xl relative"
                    style="width: {percentage}%"
                  >
                    <!-- Subtle shine effect -->
                    <div
                      class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent rounded-xl"
                    ></div>
                  </div>
                </div>
                <!-- Percentage overlay for larger bars -->
                {#if percentage > 20}
                  <div
                    class="absolute inset-0 flex items-center justify-center text-xs font-bold text-white pointer-events-none"
                    style="text-shadow: 0 1px 2px rgba(0,0,0,0.5), 0 0 4px rgba(0,0,0,0.3);"
                  >
                    {percentage.toFixed(0)}%
                  </div>
                {/if}
              </div>
              <div
                class="text-sm font-bold text-gray-900 w-20 text-right bg-white border border-blue-200 px-3 py-2 rounded-lg group-hover:bg-blue-50 transition-colors"
              >
                {count.toLocaleString()}
                <span
                  class="block text-xs text-blue-600 font-normal"
                >
                  {percentage.toFixed(1)}%
                </span>
              </div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
{:else if chart.chart_type === "choice"}
  <!-- Choice Distribution Chart -->
  {#if chart.data.is_multi_choice && chart.data.combinations && chartViewModes.get(chart.question_id) === "combinations"}
    <!-- Combination view for multi-choice -->
    <div class="space-y-3">
      {#each chart.data.combinations.slice(0, 8) as combo, index}
        <div
          class="group p-3 hover:bg-indigo-50 rounded-xl border border-indigo-100 hover:border-indigo-200 transition-all duration-200"
        >
          <div class="flex items-center gap-4">
            <div
              class="w-6 h-6 rounded-lg {getColorForChoice(
                index,
              )} flex-shrink-0 shadow-sm group-hover:scale-105 transition-transform flex items-center justify-center"
            >
              <div class="w-2 h-2 bg-white rounded-full"></div>
            </div>
            <div class="flex-1 min-w-0">
              <span
                class="text-sm font-semibold text-gray-800 block truncate"
              >
                {combo.options.join(" + ")}
              </span>
              <span class="text-xs text-indigo-600 font-medium">
                {combo.options.length} options combined
              </span>
            </div>
            <div class="flex items-center gap-4">
              <div class="relative w-28">
                <div
                  class="bg-gradient-to-r from-gray-200 to-gray-100 rounded-xl h-5 overflow-hidden shadow-inner"
                >
                  <div
                    class="h-full {getColorForChoice(
                      index,
                    )} transition-all duration-700 ease-out rounded-xl relative shadow-sm"
                    style="width: {combo.percentage}%"
                  >
                    <!-- Subtle shine effect -->
                    <div
                      class="absolute inset-0 bg-gradient-to-r from-transparent via-white/30 to-transparent rounded-xl"
                    ></div>
                  </div>
                </div>
                {#if combo.percentage > 15}
                  <div
                    class="absolute inset-0 flex items-center justify-center text-xs font-bold text-white pointer-events-none"
                    style="text-shadow: 0 1px 2px rgba(0,0,0,0.5), 0 0 4px rgba(0,0,0,0.3);"
                  >
                    {combo.percentage.toFixed(0)}%
                  </div>
                {/if}
              </div>
              <div
                class="text-sm font-bold text-gray-900 w-20 text-right bg-white border border-indigo-200 px-3 py-2 rounded-lg group-hover:bg-indigo-50 transition-colors"
              >
                {combo.count.toLocaleString()}
                <span
                  class="block text-xs text-indigo-600 font-normal"
                >
                  {combo.percentage.toFixed(1)}%
                </span>
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>
  {:else}
    <!-- Individual options view -->
    <div class="space-y-3">
      {#if chart.data.options}
        {#each Object.entries(chart.data.options).sort(([, a], [, b]) => b - a) as [option, count], index}
          {@const percentage = chart.data.percentages?.[option] || 0}
          <div
            class="group p-3 hover:bg-purple-50 rounded-xl border border-purple-100 hover:border-purple-200 transition-all duration-200"
          >
            <div class="flex items-center gap-4">
              <div
                class="w-6 h-6 rounded-lg {getColorForChoice(
                  index,
                )} flex-shrink-0 shadow-sm group-hover:scale-105 transition-transform flex items-center justify-center"
              >
                <div class="w-2 h-2 bg-white rounded-full"></div>
              </div>
              <span
                class="text-sm font-semibold text-gray-800 flex-1 truncate"
                >{option}</span
              >
              <div class="flex items-center gap-4">
                <div class="relative w-28">
                  <div
                    class="bg-gradient-to-r from-gray-200 to-gray-100 rounded-xl h-5 overflow-hidden shadow-inner"
                  >
                    <div
                      class="h-full {getColorForChoice(
                        index,
                      )} transition-all duration-700 ease-out rounded-xl relative shadow-sm"
                      style="width: {percentage}%"
                    >
                      <!-- Subtle shine effect -->
                      <div
                        class="absolute inset-0 bg-gradient-to-r from-transparent via-white/30 to-transparent rounded-xl"
                      ></div>
                    </div>
                  </div>
                  {#if percentage > 15}
                    <div
                      class="absolute inset-0 flex items-center justify-center text-xs font-bold text-white pointer-events-none"
                      style="text-shadow: 0 1px 2px rgba(0,0,0,0.5), 0 0 4px rgba(0,0,0,0.3);"
                    >
                      {percentage.toFixed(0)}%
                    </div>
                  {/if}
                </div>
                <div
                  class="text-sm font-bold text-gray-900 w-20 text-right bg-white border border-purple-200 px-3 py-2 rounded-lg group-hover:bg-purple-50 transition-colors"
                >
                  {count.toLocaleString()}
                  <span
                    class="block text-xs text-purple-600 font-normal"
                  >
                    {percentage.toFixed(1)}%
                  </span>
                </div>
              </div>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  {/if}
{:else if chart.chart_type === "text_sentiment"}
  <!-- Text Sentiment Analysis -->
  <div class="space-y-4">
    <!-- Sentiment overview -->
    <div class="grid grid-cols-3 gap-3 mb-6">
      <div
        class="text-center p-4 bg-gradient-to-br from-green-50 to-emerald-50 rounded-xl border border-green-200 relative overflow-hidden group hover:shadow-lg transition-all duration-200"
      >
        <div class="absolute top-1 right-1 text-lg opacity-20">😋</div>
        <div
          class="w-10 h-10 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform shadow-lg"
        >
          <ThumbsUp class="h-5 w-5 text-white" />
        </div>
        <div class="text-2xl font-black text-green-700 mb-1">
          {chart.data.positive || 0}
        </div>
        <div
          class="text-xs font-bold text-green-600 uppercase tracking-wide"
        >
          Loved It
        </div>
      </div>
      <div
        class="text-center p-4 bg-gradient-to-br from-amber-50 to-yellow-50 rounded-xl border border-amber-200 relative overflow-hidden group hover:shadow-lg transition-all duration-200"
      >
        <div class="absolute top-1 right-1 text-lg opacity-20">🤔</div>
        <div
          class="w-10 h-10 bg-gradient-to-br from-amber-400 to-yellow-500 rounded-xl flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform shadow-lg"
        >
          <Meh class="h-5 w-5 text-white" />
        </div>
        <div class="text-2xl font-black text-amber-700 mb-1">
          {chart.data.neutral || 0}
        </div>
        <div
          class="text-xs font-bold text-amber-600 uppercase tracking-wide"
        >
          It's OK
        </div>
      </div>
      <div
        class="text-center p-4 bg-gradient-to-br from-red-50 to-rose-50 rounded-xl border border-red-200 relative overflow-hidden group hover:shadow-lg transition-all duration-200"
      >
        <div class="absolute top-1 right-1 text-lg opacity-20">😞</div>
        <div
          class="w-10 h-10 bg-gradient-to-br from-red-500 to-rose-600 rounded-xl flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform shadow-lg"
        >
          <ThumbsDown class="h-5 w-5 text-white" />
        </div>
        <div class="text-2xl font-black text-red-700 mb-1">
          {chart.data.negative || 0}
        </div>
        <div
          class="text-xs font-bold text-red-600 uppercase tracking-wide"
        >
          Not Great
        </div>
      </div>
    </div>

    <!-- Sample responses -->
    {#if chart.data.samples && chart.data.samples.length > 0}
      <div class="bg-blue-50 rounded-lg p-4 border border-blue-100">
        <div class="flex items-center gap-2 mb-3">
          <MessageSquare class="h-4 w-4 text-blue-600" />
          <h4 class="text-sm font-semibold text-gray-900">
            Sample Responses
          </h4>
        </div>
        <div class="space-y-2">
          {#each chart.data.samples.slice(0, 3) as sample, index}
            <div
              class="bg-white p-3 rounded-lg border border-blue-100"
            >
              <p class="text-sm text-gray-700 italic">
                "{sample}"
              </p>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Keywords -->
    {#if chart.data.keywords && chart.data.keywords.length > 0}
      <div
        class="bg-purple-50 rounded-lg p-4 border border-purple-100"
      >
        <div class="flex items-center gap-2 mb-3">
          <span class="text-purple-600 text-sm font-bold">#</span>
          <h4 class="text-sm font-semibold text-gray-900">Keywords</h4>
        </div>
        <div class="flex flex-wrap gap-2">
          {#each chart.data.keywords as keyword}
            <span
              class="px-3 py-1 bg-purple-100 text-purple-800 text-xs font-medium rounded-full"
            >
              {keyword}
            </span>
          {/each}
        </div>
      </div>
    {/if}
  </div>
{/if}