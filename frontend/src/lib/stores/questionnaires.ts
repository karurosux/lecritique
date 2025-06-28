import { writable, derived } from 'svelte/store';
import { QuestionnaireApi, type Questionnaire, type Question, type GeneratedQuestion } from '$lib/api/questionnaire';

interface QuestionnaireState {
  questionnaires: Questionnaire[];
  currentQuestionnaire: Questionnaire | null;
  generatedQuestions: GeneratedQuestion[];
  loading: boolean;
  error: string | null;
}

const initialState: QuestionnaireState = {
  questionnaires: [],
  currentQuestionnaire: null,
  generatedQuestions: [],
  loading: false,
  error: null
};

function createQuestionnaireStore() {
  const { subscribe, set, update } = writable<QuestionnaireState>(initialState);

  return {
    subscribe,
    
    // Load questionnaires for a restaurant
    async loadQuestionnaires(restaurantId: string) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        const questionnaires = await QuestionnaireApi.listQuestionnaires(restaurantId);
        update(state => ({ ...state, questionnaires, loading: false }));
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
      }
    },

    // Load a specific questionnaire
    async loadQuestionnaire(restaurantId: string, questionnaireId: string) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        const questionnaire = await QuestionnaireApi.getQuestionnaire(restaurantId, questionnaireId);
        update(state => ({ ...state, currentQuestionnaire: questionnaire, loading: false }));
        return questionnaire;
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Create new questionnaire
    async createQuestionnaire(restaurantId: string, data: any) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        const questionnaire = await QuestionnaireApi.createQuestionnaire(restaurantId, data);
        update(state => ({ 
          ...state, 
          questionnaires: [...state.questionnaires, questionnaire],
          currentQuestionnaire: questionnaire,
          loading: false 
        }));
        return questionnaire;
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Update questionnaire
    async updateQuestionnaire(restaurantId: string, questionnaireId: string, data: any) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        const updatedQuestionnaire = await QuestionnaireApi.updateQuestionnaire(restaurantId, questionnaireId, data);
        update(state => ({
          ...state,
          questionnaires: state.questionnaires.map(q => 
            q.id === questionnaireId ? updatedQuestionnaire : q
          ),
          currentQuestionnaire: state.currentQuestionnaire?.id === questionnaireId 
            ? updatedQuestionnaire 
            : state.currentQuestionnaire,
          loading: false
        }));
        return updatedQuestionnaire;
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Delete questionnaire
    async deleteQuestionnaire(restaurantId: string, questionnaireId: string) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        await QuestionnaireApi.deleteQuestionnaire(restaurantId, questionnaireId);
        update(state => ({
          ...state,
          questionnaires: state.questionnaires.filter(q => q.id !== questionnaireId),
          currentQuestionnaire: state.currentQuestionnaire?.id === questionnaireId 
            ? null 
            : state.currentQuestionnaire,
          loading: false
        }));
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Add question to current questionnaire
    async addQuestion(restaurantId: string, questionnaireId: string, question: any) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        const newQuestion = await QuestionnaireApi.addQuestion(restaurantId, questionnaireId, question);
        update(state => {
          const updatedQuestionnaire = state.currentQuestionnaire?.id === questionnaireId
            ? {
                ...state.currentQuestionnaire,
                questions: [...(state.currentQuestionnaire.questions || []), newQuestion]
              }
            : state.currentQuestionnaire;

          return {
            ...state,
            currentQuestionnaire: updatedQuestionnaire,
            loading: false
          };
        });
        return newQuestion;
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Update question
    async updateQuestion(restaurantId: string, questionnaireId: string, questionId: string, question: any) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        const updatedQuestion = await QuestionnaireApi.updateQuestion(restaurantId, questionnaireId, questionId, question);
        update(state => {
          const updatedQuestionnaire = state.currentQuestionnaire?.id === questionnaireId
            ? {
                ...state.currentQuestionnaire,
                questions: state.currentQuestionnaire.questions?.map(q => 
                  q.id === questionId ? updatedQuestion : q
                ) || []
              }
            : state.currentQuestionnaire;

          return {
            ...state,
            currentQuestionnaire: updatedQuestionnaire,
            loading: false
          };
        });
        return updatedQuestion;
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Delete question
    async deleteQuestion(restaurantId: string, questionnaireId: string, questionId: string) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        await QuestionnaireApi.deleteQuestion(restaurantId, questionnaireId, questionId);
        update(state => {
          const updatedQuestionnaire = state.currentQuestionnaire?.id === questionnaireId
            ? {
                ...state.currentQuestionnaire,
                questions: state.currentQuestionnaire.questions?.filter(q => q.id !== questionId) || []
              }
            : state.currentQuestionnaire;

          return {
            ...state,
            currentQuestionnaire: updatedQuestionnaire,
            loading: false
          };
        });
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Reorder questions
    async reorderQuestions(restaurantId: string, questionnaireId: string, questionIds: string[]) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        await QuestionnaireApi.reorderQuestions(restaurantId, questionnaireId, questionIds);
        // Reload the questionnaire to get the updated order
        const questionnaire = await QuestionnaireApi.getQuestionnaire(restaurantId, questionnaireId);
        update(state => ({ ...state, currentQuestionnaire: questionnaire, loading: false }));
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Generate AI questions
    async generateQuestions(dishId: string) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        const questions = await QuestionnaireApi.generateQuestions(dishId);
        update(state => ({ ...state, generatedQuestions: questions, loading: false }));
        return questions;
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Generate and save questionnaire
    async generateAndSaveQuestionnaire(dishId: string, data: any) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        const questionnaire = await QuestionnaireApi.generateAndSaveQuestionnaire(dishId, data);
        update(state => ({ 
          ...state, 
          questionnaires: [...state.questionnaires, questionnaire],
          currentQuestionnaire: questionnaire,
          loading: false 
        }));
        return questionnaire;
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Clear current questionnaire
    clearCurrent() {
      update(state => ({ ...state, currentQuestionnaire: null }));
    },

    // Clear generated questions
    clearGenerated() {
      update(state => ({ ...state, generatedQuestions: [] }));
    },

    // Clear error
    clearError() {
      update(state => ({ ...state, error: null }));
    },

    // Reset store
    reset() {
      set(initialState);
    }
  };
}

export const questionnaireStore = createQuestionnaireStore();

// Derived stores for easy access
export const questionnaires = derived(questionnaireStore, $state => $state.questionnaires);
export const currentQuestionnaire = derived(questionnaireStore, $state => $state.currentQuestionnaire);
export const generatedQuestions = derived(questionnaireStore, $state => $state.generatedQuestions);
export const questionnaireLoading = derived(questionnaireStore, $state => $state.loading);
export const questionnaireError = derived(questionnaireStore, $state => $state.error);