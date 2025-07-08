import { getApiClient, handleApiError } from './client';
import type { 
  ModelsQuestion,
  ModelsCreateQuestionRequest,
  ModelsUpdateQuestionRequest,
  ModelsQuestionType
} from './api';

// Re-export types for easy access
export type Question = ModelsQuestion;
export type CreateQuestionRequest = ModelsCreateQuestionRequest;
export type UpdateQuestionRequest = ModelsUpdateQuestionRequest;
export type QuestionType = ModelsQuestionType;

// Simplified Question API for dish-based questions
export class QuestionApi {
  
  // Get all questions for a dish
  static async getQuestionsByDish(restaurantId: string, dishId: string): Promise<Question[]> {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsDishesQuestionsList(restaurantId, dishId);
      return response.data.data || [];
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Create a new question for a dish
  static async createQuestion(restaurantId: string, dishId: string, data: CreateQuestionRequest): Promise<Question> {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsDishesQuestionsCreate(restaurantId, dishId, data);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Get a specific question
  static async getQuestion(restaurantId: string, dishId: string, questionId: string): Promise<Question> {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsDishesQuestionsDetail(restaurantId, dishId, questionId);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Update an existing question
  static async updateQuestion(restaurantId: string, dishId: string, questionId: string, data: UpdateQuestionRequest): Promise<Question> {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsDishesQuestionsUpdate(restaurantId, dishId, questionId, data);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Delete a question
  static async deleteQuestion(restaurantId: string, dishId: string, questionId: string): Promise<void> {
    try {
      const api = getApiClient();
      await api.api.v1RestaurantsDishesQuestionsDelete(restaurantId, dishId, questionId);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Reorder questions for a dish
  static async reorderQuestions(restaurantId: string, dishId: string, questionIds: string[]): Promise<void> {
    try {
      const api = getApiClient();
      await api.api.v1RestaurantsDishesQuestionsReorderCreate(restaurantId, dishId, questionIds);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Get questions for a dish (public endpoint for customer feedback)
  static async getPublicQuestions(restaurantId: string, dishId: string): Promise<{ dish: any; questions: Question[] }> {
    try {
      const api = getApiClient();
      const response = await api.restaurant.dishesQuestionsList(restaurantId, dishId);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Generate AI questions for a dish (preview only)
  static async generateQuestions(dishId: string): Promise<any[]> {
    try {
      const api = getApiClient();
      const response = await api.api.v1AiGenerateQuestionsCreate(dishId);
      return response.data.data || [];
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Generate and save complete questionnaire for a dish
  static async generateAndSaveQuestionnaire(dishId: string, data: any): Promise<any> {
    try {
      const api = getApiClient();
      const response = await api.api.v1AiGenerateQuestionnaireCreate(dishId, data);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }
}