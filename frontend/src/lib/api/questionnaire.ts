import { getApiClient, handleApiError } from './client';
import type { 
  ModelsQuestionnaire,
  ModelsCreateQuestionnaireRequest,
  ModelsGenerateQuestionnaireRequest,
  ModelsGeneratedQuestion,
  ModelsQuestionType,
  ModelsQuestion
} from './api';

export type Questionnaire = ModelsQuestionnaire;
export type CreateQuestionnaireRequest = ModelsCreateQuestionnaireRequest;
export type GenerateQuestionnaireRequest = ModelsGenerateQuestionnaireRequest;
export type GeneratedQuestion = ModelsGeneratedQuestion;
export type QuestionType = ModelsQuestionType;
export type Question = ModelsQuestion;

export class QuestionnaireApi {
  
  static async listQuestionnaires(organizationId: string): Promise<Questionnaire[]> {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsQuestionnairesList(organizationId);
      return response.data.data || [];
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  static async createQuestionnaire(organizationId: string, data: CreateQuestionnaireRequest): Promise<Questionnaire> {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsQuestionnairesCreate(organizationId, data);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  static async getQuestionnaire(organizationId: string, questionnaireId: string): Promise<Questionnaire> {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsQuestionnairesDetail(organizationId, questionnaireId);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  static async updateQuestionnaire(organizationId: string, questionnaireId: string, data: Partial<Questionnaire>): Promise<Questionnaire> {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsQuestionnairesUpdate(organizationId, questionnaireId, data as ModelsQuestionnaire);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  static async deleteQuestionnaire(organizationId: string, questionnaireId: string): Promise<void> {
    try {
      const api = getApiClient();
      await api.api.v1OrganizationsQuestionnairesDelete(organizationId, questionnaireId);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  static async generateQuestions(productId: string): Promise<GeneratedQuestion[]> {
    try {
      const api = getApiClient();
      const response = await api.api.v1AiGenerateQuestionsCreate(productId);
      return response.data.data || [];
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  static async generateAndSaveQuestionnaire(productId: string, data: GenerateQuestionnaireRequest): Promise<Questionnaire> {
    try {
      const api = getApiClient();
      const response = await api.api.v1AiGenerateQuestionnaireCreate(productId, data);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  static async addQuestion(organizationId: string, questionnaireId: string, question: Question): Promise<Question> {
    try {
      const api = getApiClient();
      const response = await api.organizations.questionnairesQuestionsCreate(questionnaireId, organizationId, question);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Update a question
  static async updateQuestion(organizationId: string, questionnaireId: string, questionId: string, question: Question): Promise<Question> {
    try {
      const api = getApiClient();
      const response = await api.organizations.questionnairesQuestionsUpdate(questionnaireId, questionId, organizationId, question);
      return response.data.data!;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Delete a question
  static async deleteQuestion(organizationId: string, questionnaireId: string, questionId: string): Promise<void> {
    try {
      const api = getApiClient();
      await api.organizations.questionnairesQuestionsDelete(questionnaireId, questionId, organizationId);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }

  // Reorder questions
  static async reorderQuestions(organizationId: string, questionnaireId: string, questionIds: string[]): Promise<void> {
    try {
      const api = getApiClient();
      await api.organizations.questionnairesReorderCreate(questionnaireId, organizationId, questionIds);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }
}
