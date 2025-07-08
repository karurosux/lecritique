import { writable, derived } from 'svelte/store';
import { QuestionApi, type Question, type CreateQuestionRequest, type UpdateQuestionRequest } from '$lib/api/question';

interface QuestionState {
  questions: Question[];
  loading: boolean;
  error: string | null;
}

const initialState: QuestionState = {
  questions: [],
  loading: false,
  error: null
};

function createQuestionStore() {
  const { subscribe, set, update } = writable<QuestionState>(initialState);

  return {
    subscribe,
    
    // Load questions for a dish
    async loadQuestions(restaurantId: string, dishId: string) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        const questions = await QuestionApi.getQuestionsByDish(restaurantId, dishId);
        update(state => ({ ...state, questions, loading: false }));
        return questions;
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Create new question
    async createQuestion(restaurantId: string, dishId: string, data: CreateQuestionRequest) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        const question = await QuestionApi.createQuestion(restaurantId, dishId, data);
        update(state => ({ 
          ...state, 
          questions: [...state.questions, question].sort((a, b) => (a.display_order || 0) - (b.display_order || 0)),
          loading: false 
        }));
        return question;
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Update question
    async updateQuestion(restaurantId: string, dishId: string, questionId: string, data: UpdateQuestionRequest) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        const updatedQuestion = await QuestionApi.updateQuestion(restaurantId, dishId, questionId, data);
        update(state => ({
          ...state,
          questions: state.questions.map(q => 
            q.id === questionId ? updatedQuestion : q
          ),
          loading: false
        }));
        return updatedQuestion;
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Delete question
    async deleteQuestion(restaurantId: string, dishId: string, questionId: string) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        await QuestionApi.deleteQuestion(restaurantId, dishId, questionId);
        update(state => ({
          ...state,
          questions: state.questions.filter(q => q.id !== questionId),
          loading: false
        }));
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Reorder questions
    async reorderQuestions(restaurantId: string, dishId: string, questionIds: string[]) {
      update(state => ({ ...state, loading: true, error: null }));
      
      try {
        await QuestionApi.reorderQuestions(restaurantId, dishId, questionIds);
        // Reload questions to get updated order
        const questions = await QuestionApi.getQuestionsByDish(restaurantId, dishId);
        update(state => ({ ...state, questions, loading: false }));
      } catch (error) {
        update(state => ({ ...state, loading: false, error: error.message }));
        throw error;
      }
    },

    // Clear questions
    clear() {
      set(initialState);
    },

    // Clear error
    clearError() {
      update(state => ({ ...state, error: null }));
    }
  };
}

export const questionStore = createQuestionStore();

// Derived stores for easy access
export const questions = derived(questionStore, $state => $state.questions);
export const questionsLoading = derived(questionStore, $state => $state.loading);
export const questionsError = derived(questionStore, $state => $state.error);