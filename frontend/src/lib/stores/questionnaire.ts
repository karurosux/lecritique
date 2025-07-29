import { writable, derived } from 'svelte/store';
import { getApiClient, handleApiError } from '$lib/api/client';

export type QuestionType =
  | 'rating'
  | 'scale'
  | 'single_choice'
  | 'multiple_choice'
  | 'yes_no'
  | 'text';

export interface QuestionOption {
  id: string;
  text: string;
  value?: string;
}

export interface Question {
  id: string;
  text: string;
  type: QuestionType;
  required: boolean;
  options?: QuestionOption[];
  min_value?: number;
  max_value?: number;
  min_label?: string;
  max_label?: string;
  order: number;
}

export interface Questionnaire {
  id: string;
  organization_id: string;
  location_id?: string;
  name: string;
  description?: string;
  questions: Question[];
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
}

interface QuestionnaireState {
  questionnaire: Questionnaire | null;
  loading: boolean;
  error: string | null;
  cache: Map<string, { data: Questionnaire; timestamp: number }>;
}

const CACHE_DURATION = 5 * 60 * 1000; // 5 minutes

function createQuestionnaireStore() {
  const { subscribe, set, update } = writable<QuestionnaireState>({
    questionnaire: null,
    loading: false,
    error: null,
    cache: new Map(),
  });

  async function fetchQuestionnaire(
    organizationId: string,
    locationId?: string
  ) {
    const cacheKey = `${organizationId}-${locationId || 'default'}`;

    update(state => {
      // Check cache
      const cached = state.cache.get(cacheKey);
      if (cached && Date.now() - cached.timestamp < CACHE_DURATION) {
        return {
          ...state,
          questionnaire: cached.data,
          loading: false,
          error: null,
        };
      }

      return { ...state, loading: true, error: null };
    });

    try {
      const api = getApiClient();

      // For now, we'll use a fallback since the API might not have location-based questionnaires
      // We'll try to get the questionnaire for the organization
      const response = await api.api.v1PublicQuestionnaireDetail(
        organizationId,
        'default'
      );

      // If the API returns an empty response or error, use default questions
      const questionnaire: Questionnaire =
        response.data &&
        typeof response.data === 'object' &&
        'questions' in response.data
          ? (response.data as Questionnaire)
          : {
              id: 'default',
              organization_id: organizationId,
              location_id: locationId,
              name: 'Customer Feedback',
              description: 'Help us improve your dining experience',
              is_active: true,
              questions: getDefaultQuestions(),
            };

      update(state => {
        // Update cache
        state.cache.set(cacheKey, {
          data: questionnaire,
          timestamp: Date.now(),
        });

        return {
          ...state,
          questionnaire,
          loading: false,
          error: null,
        };
      });

      return questionnaire;
    } catch (err) {
      const errorMessage = handleApiError(err);

      // If API fails, use default questions
      const fallbackQuestionnaire: Questionnaire = {
        id: 'default',
        organization_id: organizationId,
        location_id: locationId,
        name: 'Customer Feedback',
        description: 'Help us improve your dining experience',
        is_active: true,
        questions: getDefaultQuestions(),
      };

      update(state => ({
        ...state,
        questionnaire: fallbackQuestionnaire,
        loading: false,
        error: errorMessage,
      }));

      return fallbackQuestionnaire;
    }
  }

  async function fetchProductQuestions(
    organizationId: string,
    productId: string
  ) {
    const cacheKey = `${organizationId}-product-${productId}`;

    update(state => {
      // Check cache
      const cached = state.cache.get(cacheKey);
      if (cached && Date.now() - cached.timestamp < CACHE_DURATION) {
        return {
          ...state,
          questionnaire: cached.data,
          loading: false,
          error: null,
        };
      }

      return { ...state, loading: true, error: null };
    });

    try {
      const api = getApiClient();

      // Fetch product questions from the public API
      const response =
        await api.api.v1PublicOrganizationProductsQuestionsDetail(
          organizationId,
          productId
        );

      if (response.data && response.data.success && response.data.data) {
        const { product, questions } = response.data.data;

        // Transform product questions into questionnaire format
        const questionnaire: Questionnaire = {
          id: `product-${productId}`,
          organization_id: organizationId,
          name: `${product.name} Feedback`,
          description: `Tell us about your experience with ${product.name}`,
          is_active: true,
          questions: questions.map((q: any, index: number) => ({
            id: q.id,
            text: q.text,
            type: q.type,
            required: q.is_required,
            options: q.options || [],
            min_value: q.min_value,
            max_value: q.max_value,
            min_label: q.min_label,
            max_label: q.max_label,
            order: q.display_order || index + 1,
          })),
        };

        update(state => {
          // Update cache
          state.cache.set(cacheKey, {
            data: questionnaire,
            timestamp: Date.now(),
          });

          return {
            ...state,
            questionnaire,
            loading: false,
            error: null,
          };
        });

        return questionnaire;
      } else {
        throw new Error('No questions found for this product');
      }
    } catch (err) {
      const errorMessage = handleApiError(err);

      // If API fails, use default questions
      const fallbackQuestionnaire: Questionnaire = {
        id: 'default',
        organization_id: organizationId,
        name: 'Customer Feedback',
        description: 'Help us improve your dining experience',
        is_active: true,
        questions: getDefaultQuestions(),
      };

      update(state => ({
        ...state,
        questionnaire: fallbackQuestionnaire,
        loading: false,
        error: errorMessage,
      }));

      return fallbackQuestionnaire;
    }
  }

  function clearCache() {
    update(state => ({
      ...state,
      cache: new Map(),
    }));
  }

  function reset() {
    set({
      questionnaire: null,
      loading: false,
      error: null,
      cache: new Map(),
    });
  }

  return {
    subscribe,
    fetchQuestionnaire,
    fetchProductQuestions,
    clearCache,
    reset,
  };
}

// Default questions as fallback
function getDefaultQuestions(): Question[] {
  return [
    {
      id: 'q1',
      text: 'How would you rate the food quality?',
      type: 'rating',
      required: true,
      min_value: 1,
      max_value: 5,
      order: 1,
    },
    {
      id: 'q2',
      text: 'How would you rate the service?',
      type: 'rating',
      required: true,
      min_value: 1,
      max_value: 5,
      order: 2,
    },
    {
      id: 'q3',
      text: 'How would you rate the ambiance?',
      type: 'rating',
      required: true,
      min_value: 1,
      max_value: 5,
      order: 3,
    },
    {
      id: 'q4',
      text: 'How would you rate the value for money?',
      type: 'rating',
      required: true,
      min_value: 1,
      max_value: 5,
      order: 4,
    },
    {
      id: 'q5',
      text: 'Was the food served at the right temperature?',
      type: 'single_choice',
      required: true,
      options: [
        { id: 'opt1', text: 'Too cold' },
        { id: 'opt2', text: 'Just right' },
        { id: 'opt3', text: 'Too hot' },
      ],
      order: 5,
    },
    {
      id: 'q6',
      text: 'How was the waiting time?',
      type: 'single_choice',
      required: true,
      options: [
        { id: 'opt4', text: 'Too long' },
        { id: 'opt5', text: 'Reasonable' },
        { id: 'opt6', text: 'Very quick' },
      ],
      order: 6,
    },
    {
      id: 'q7',
      text: 'Would you recommend us to a friend?',
      type: 'yes_no',
      required: true,
      order: 7,
    },
    {
      id: 'q8',
      text: 'Any additional comments or suggestions?',
      type: 'text',
      required: false,
      order: 8,
    },
  ];
}

export const questionnaireStore = createQuestionnaireStore();

// Derived store for easy access to current questionnaire
export const currentQuestionnaire = derived(
  questionnaireStore,
  $store => $store.questionnaire
);

// Derived store for loading state
export const questionnaireLoading = derived(
  questionnaireStore,
  $store => $store.loading
);

// Derived store for error state
export const questionnaireError = derived(
  questionnaireStore,
  $store => $store.error
);
