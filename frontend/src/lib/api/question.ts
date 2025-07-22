import { getApiClient, handleApiError, getAuthToken } from './client';
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

// Simplified Question API for product-based questions
export class QuestionApi {
  
  // Get all questions for a product
  static async getQuestionsByProduct(organizationId: string, productId: string): Promise<Question[]> {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsProductesQuestionsList(organizationId, productId);
      return response.data.data || [];
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Create a new question for a product
  static async createQuestion(organizationId: string, productId: string, data: CreateQuestionRequest): Promise<Question> {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsProductesQuestionsCreate(organizationId, productId, data);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Get a specific question
  static async getQuestion(organizationId: string, productId: string, questionId: string): Promise<Question> {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsProductesQuestionsDetail(organizationId, productId, questionId);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Update an existing question
  static async updateQuestion(organizationId: string, productId: string, questionId: string, data: UpdateQuestionRequest): Promise<Question> {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsProductesQuestionsUpdate(organizationId, productId, questionId, data);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Delete a question
  static async deleteQuestion(organizationId: string, productId: string, questionId: string): Promise<void> {
    try {
      const api = getApiClient();
      await api.api.v1OrganizationsProductesQuestionsDelete(organizationId, productId, questionId);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Reorder questions for a product
  static async reorderQuestions(organizationId: string, productId: string, questionIds: string[]): Promise<void> {
    try {
      const api = getApiClient();
      await api.api.v1OrganizationsProductesQuestionsReorderCreate(organizationId, productId, questionIds);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Get questions for a product (public endpoint for customer feedback)
  static async getPublicQuestions(organizationId: string, productId: string): Promise<{ product: any; questions: Question[] }> {
    try {
      const api = getApiClient();
      const response = await api.organization.productsQuestionsList(organizationId, productId);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Generate AI questions for a product (preview only)
  static async generateQuestions(organizationId: string, productId: string): Promise<any[]> {
    try {
      const token = getAuthToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const response = await fetch(`http://localhost:8080/api/v1/organizations/${organizationId}/products/${productId}/ai/generate-questions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();
      return data?.data || [];
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Generate and save complete questionnaire for a product
  static async generateAndSaveQuestionnaire(productId: string, data: any): Promise<any> {
    try {
      const api = getApiClient();
      const response = await api.api.v1AiGenerateQuestionnaireCreate(productId, data);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }
}
