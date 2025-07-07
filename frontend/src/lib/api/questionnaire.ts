import { getApiClient, handleApiError } from './client';
import type { 
  ModelsQuestionnaire,
  ModelsCreateQuestionnaireRequest,
  ModelsGenerateQuestionnaireRequest,
  ModelsGeneratedQuestion,
  ModelsQuestionType,
  ModelsQuestion
} from './api';

// Re-export types for easy access
export type Questionnaire = ModelsQuestionnaire;
export type CreateQuestionnaireRequest = ModelsCreateQuestionnaireRequest;
export type GenerateQuestionnaireRequest = ModelsGenerateQuestionnaireRequest;
export type GeneratedQuestion = ModelsGeneratedQuestion;
export type QuestionType = ModelsQuestionType;
export type Question = ModelsQuestion;

// Use the generated API client
export class QuestionnaireApi {
  
  // List all questionnaires for a restaurant
  static async listQuestionnaires(restaurantId: string): Promise<Questionnaire[]> {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsQuestionnairesList(restaurantId);
      return response.data.data || [];
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Create a new questionnaire
  static async createQuestionnaire(restaurantId: string, data: CreateQuestionnaireRequest): Promise<Questionnaire> {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsQuestionnairesCreate(restaurantId, data);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Get questionnaire by ID
  static async getQuestionnaire(restaurantId: string, questionnaireId: string): Promise<Questionnaire> {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsQuestionnairesDetail(restaurantId, questionnaireId);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Update questionnaire
  static async updateQuestionnaire(restaurantId: string, questionnaireId: string, data: Partial<Questionnaire>): Promise<Questionnaire> {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsQuestionnairesUpdate(restaurantId, questionnaireId, data as ModelsQuestionnaire);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Delete questionnaire
  static async deleteQuestionnaire(restaurantId: string, questionnaireId: string): Promise<void> {
    try {
      const api = getApiClient();
      await api.api.v1RestaurantsQuestionnairesDelete(restaurantId, questionnaireId);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Generate AI questions for a dish (preview only)
  static async generateQuestions(dishId: string): Promise<GeneratedQuestion[]> {
    try {
      const api = getApiClient();
      const response = await api.api.v1AiGenerateQuestionsCreate(dishId);
      return response.data.data || [];
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Generate and save complete questionnaire for a dish
  static async generateAndSaveQuestionnaire(dishId: string, data: GenerateQuestionnaireRequest): Promise<Questionnaire> {
    try {
      const api = getApiClient();
      const response = await api.api.v1AiGenerateQuestionnaireCreate(dishId, data);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Add a question to a questionnaire
  static async addQuestion(restaurantId: string, questionnaireId: string, question: Question): Promise<Question> {
    try {
      const api = getApiClient();
      const response = await api.restaurants.questionnairesQuestionsCreate(questionnaireId, restaurantId, question);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Update a question
  static async updateQuestion(restaurantId: string, questionnaireId: string, questionId: string, question: Question): Promise<Question> {
    try {
      const api = getApiClient();
      const response = await api.restaurants.questionnairesQuestionsUpdate(questionnaireId, questionId, restaurantId, question);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Delete a question
  static async deleteQuestion(restaurantId: string, questionnaireId: string, questionId: string): Promise<void> {
    try {
      const api = getApiClient();
      await api.restaurants.questionnairesQuestionsDelete(questionnaireId, questionId, restaurantId);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Reorder questions
  static async reorderQuestions(restaurantId: string, questionnaireId: string, questionIds: string[]): Promise<void> {
    try {
      const api = getApiClient();
      await api.restaurants.questionnairesReorderCreate(questionnaireId, restaurantId, questionIds);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }
}