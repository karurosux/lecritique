/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export enum ModelsSubscriptionStatus {
  SubscriptionActive = "SubscriptionActive",
  SubscriptionPending = "SubscriptionPending",
  SubscriptionCanceled = "SubscriptionCanceled",
  SubscriptionExpired = "SubscriptionExpired",
}

export enum ModelsQuestionType {
  QuestionTypeRating = "QuestionTypeRating",
  QuestionTypeScale = "QuestionTypeScale",
  QuestionTypeMultiChoice = "QuestionTypeMultiChoice",
  QuestionTypeSingleChoice = "QuestionTypeSingleChoice",
  QuestionTypeText = "QuestionTypeText",
  QuestionTypeYesNo = "QuestionTypeYesNo",
}

export enum ModelsQRCodeType {
  QRCodeTypeTable = "QRCodeTypeTable",
  QRCodeTypeLocation = "QRCodeTypeLocation",
  QRCodeTypeTakeaway = "QRCodeTypeTakeaway",
  QRCodeTypeDelivery = "QRCodeTypeDelivery",
  QRCodeTypeGeneral = "QRCodeTypeGeneral",
}

export enum ModelsMemberRole {
  RoleOwner = "RoleOwner",
  RoleAdmin = "RoleAdmin",
  RoleManager = "RoleManager",
  RoleViewer = "RoleViewer",
}

export interface HandlersAcceptInviteRequest {
  token: string;
}

export interface HandlersAuthResponse {
  account?: any;
  token?: string;
}

export interface HandlersCardDetailsResponse {
  brand?: string;
  exp_month?: number;
  exp_year?: number;
  last4?: string;
}

export interface HandlersChangeEmailRequest {
  new_email: string;
}

export interface HandlersCheckoutResponse {
  checkout_url?: string;
  session_id?: string;
}

export interface HandlersCompleteCheckoutRequest {
  session_id: string;
}

export interface HandlersConfirmEmailChangeRequest {
  token: string;
}

export interface HandlersCreateCheckoutRequest {
  plan_id: string;
}

export interface HandlersCreateOrganizationRequest {
  description?: string;
  email?: string;
  name: string;
  phone?: string;
  website?: string;
}

export interface HandlersCreateProductRequest {
  category?: string;
  currency?: string;
  description?: string;
  name: string;
  organization_id: string;
  /** @min 0 */
  price?: number;
}

export interface HandlersCreateSubscriptionRequest {
  plan_id: string;
}

export interface HandlersGenerateQRCodeRequest {
  /**
   * @minLength 1
   * @maxLength 100
   */
  label: string;
  /** @maxLength 200 */
  location?: string;
  organization_id: string;
  type: "table" | "location" | "takeaway" | "delivery" | "general";
}

export interface HandlersGenerateQRCodeResponse {
  data?: ModelsQRCode;
  success?: boolean;
}

export interface HandlersInviteMemberRequest {
  email: string;
  role: "ADMIN" | "MANAGER" | "VIEWER";
}

export interface HandlersInvoiceResponse {
  amount_due?: number;
  amount_paid?: number;
  created_at?: string;
  currency?: string;
  hosted_invoice_url?: string;
  id?: string;
  invoice_pdf?: string;
  number?: string;
  paid_at?: string;
  status?: string;
}

export interface HandlersLoginRequest {
  email: string;
  password: string;
}

export interface HandlersPasswordResetRequest {
  email: string;
}

export interface HandlersPaymentMethodResponse {
  card?: HandlersCardDetailsResponse;
  id?: string;
  is_default?: boolean;
  type?: string;
}

export interface HandlersPortalResponse {
  portal_url?: string;
}

export interface HandlersQRCodeListResponse {
  data?: ModelsQRCode[];
  success?: boolean;
}

export interface HandlersRegisterRequest {
  email: string;
  first_name?: string;
  /** Optional invitation token */
  invitation_token?: string;
  last_name?: string;
  name?: string;
  /** @minLength 8 */
  password: string;
}

export interface HandlersResendVerificationRequest {
  email: string;
}

export interface HandlersResetPasswordRequest {
  /** @minLength 8 */
  new_password: string;
  token: string;
}

export interface HandlersSetDefaultPaymentRequest {
  payment_method_id: string;
}

export interface HandlersUpdateProfileRequest {
  /** @minLength 1 */
  name?: string;
  phone?: string;
}

export interface HandlersUpdateQRCodeRequest {
  is_active?: boolean;
  /**
   * @minLength 1
   * @maxLength 100
   */
  label?: string;
  /** @maxLength 200 */
  location?: string;
}

export interface HandlersUpdateQRCodeResponse {
  data?: ModelsQRCode;
  success?: boolean;
}

export interface HandlersUpdateRoleRequest {
  role: "ADMIN" | "MANAGER" | "VIEWER";
}

export type KyooarInternalMenuModelsProduct = object;

export interface ModelsAccount {
  created_at?: string;
  deactivation_requested_at?: string;
  email?: string;
  email_verified?: boolean;
  email_verified_at?: string;
  first_name?: string;
  id?: string;
  is_active?: boolean;
  last_name?: string;
  name?: string;
  phone?: string;
  /** Populated when needed */
  subscription?: any;
  subscription_id?: string;
  /** Organizations      []Organization  `json:"organizations,omitempty"`  // TODO: Add when organization domain is ready */
  team_members?: ModelsTeamMember[];
  updated_at?: string;
}

export interface ModelsCreateQuestionRequest {
  is_required?: boolean;
  max_label?: string;
  max_value?: number;
  min_label?: string;
  min_value?: number;
  options?: string[];
  text: string;
  type: ModelsQuestionType;
}

export interface ModelsCreateQuestionnaireRequest {
  description?: string;
  is_default?: boolean;
  name: string;
  product_id?: string;
}

export type ModelsFeedback = object;

export interface ModelsGenerateQuestionnaireRequest {
  description?: string;
  is_default?: boolean;
  name: string;
}

export interface ModelsGeneratedQuestion {
  max_label?: string;
  max_value?: number;
  min_label?: string;
  min_value?: number;
  options?: string[];
  text?: string;
  type?: ModelsQuestionType;
}

export interface ModelsLocation {
  address?: string;
  city?: string;
  country?: string;
  created_at?: string;
  id?: string;
  is_active?: boolean;
  latitude?: number;
  longitude?: number;
  name?: string;
  organization?: ModelsOrganization;
  organization_id?: string;
  postal_code?: string;
  state?: string;
  updated_at?: string;
}

export interface ModelsOrganization {
  account_id?: string;
  created_at?: string;
  description?: string;
  email?: string;
  id?: string;
  is_active?: boolean;
  locations?: ModelsLocation[];
  logo?: string;
  /** Account     Account        `json:"account,omitempty"` // TODO: Add when cross-domain refs are ready */
  name?: string;
  phone?: string;
  settings?: ModelsSettings;
  updated_at?: string;
  website?: string;
}

export interface ModelsQRCode {
  code?: string;
  created_at?: string;
  expires_at?: string;
  id?: string;
  is_active?: boolean;
  /** e.g., "Table 1", "Entrance", etc. */
  label?: string;
  last_scanned_at?: string;
  /** Free text location description */
  location?: string;
  organization?: ModelsOrganization;
  organization_id?: string;
  scans_count?: number;
  type?: ModelsQRCodeType;
  updated_at?: string;
}

export interface ModelsQuestion {
  created_at?: string;
  display_order?: number;
  id?: string;
  is_required?: boolean;
  max_label?: string;
  max_value?: number;
  min_label?: string;
  min_value?: number;
  options?: string[];
  product?: KyooarInternalMenuModelsProduct;
  product_id?: string;
  text?: string;
  type?: ModelsQuestionType;
  updated_at?: string;
}

export interface ModelsQuestionnaire {
  created_at?: string;
  description?: string;
  id?: string;
  is_active?: boolean;
  is_default?: boolean;
  name?: string;
  organization?: ModelsOrganization;
  organization_id?: string;
  product?: KyooarInternalMenuModelsProduct;
  product_id?: string;
  questions?: ModelsQuestion[];
  updated_at?: string;
}

export interface ModelsSettings {
  feedback_notification?: boolean;
  language?: string;
  low_rating_threshold?: number;
  timezone?: string;
}

export interface ModelsSubscription {
  account?: ModelsAccount;
  account_id?: string;
  cancel_at?: string;
  cancelled_at?: string;
  created_at?: string;
  current_period_end?: string;
  current_period_start?: string;
  id?: string;
  plan?: ModelsSubscriptionPlan;
  plan_id?: string;
  status?: ModelsSubscriptionStatus;
  updated_at?: string;
}

export interface ModelsSubscriptionPlan {
  code?: string;
  created_at?: string;
  currency?: string;
  description?: string;
  has_advanced_analytics?: boolean;
  /** Feature flags (as columns) */
  has_basic_analytics?: boolean;
  has_custom_branding?: boolean;
  has_feedback_explorer?: boolean;
  has_priority_support?: boolean;
  id?: string;
  interval?: string;
  is_active?: boolean;
  is_visible?: boolean;
  max_feedbacks_per_month?: number;
  /** Limits (as columns) */
  max_organizations?: number;
  max_qr_codes?: number;
  max_team_members?: number;
  name?: string;
  price?: number;
  trial_days?: number;
  updated_at?: string;
}

export interface ModelsSubscriptionUsage {
  created_at?: string;
  feedbacks_count?: number;
  id?: string;
  last_updated_at?: string;
  locations_count?: number;
  organizations_count?: number;
  period_end?: string;
  period_start?: string;
  qr_codes_count?: number;
  subscription?: ModelsSubscription;
  subscription_id?: string;
  team_members_count?: number;
  updated_at?: string;
}

export interface ModelsTeamMember {
  accepted_at?: string;
  account?: ModelsAccount;
  /** The organization account */
  account_id?: string;
  created_at?: string;
  id?: string;
  invited_at?: string;
  invited_by?: string;
  member?: ModelsAccount;
  /** The member's account ID */
  member_id?: string;
  role?: ModelsMemberRole;
  updated_at?: string;
}

export interface ModelsUpdateQuestionRequest {
  is_required?: boolean;
  max_label?: string;
  max_value?: number;
  min_label?: string;
  min_value?: number;
  options?: string[];
  text?: string;
  type?: ModelsQuestionType;
}

export interface ResponseErrorData {
  code?: string;
  details?: any;
  message?: string;
}

export interface ResponseMeta {
  pagination?: ResponsePagination;
  request_id?: string;
  timestamp?: string;
  version?: string;
}

export interface ResponsePagination {
  limit?: number;
  page?: number;
  pages?: number;
  total?: number;
}

export interface ResponseResponse {
  data?: any;
  error?: ResponseErrorData;
  meta?: ResponseMeta;
  success?: boolean;
}

export interface ServicesPermissionResponse {
  can_create?: boolean;
  current_count?: number;
  max_allowed?: number;
  reason?: string;
  subscription_status?: string;
}

import type {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  HeadersDefaults,
  ResponseType,
} from "axios";
import axios from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams
  extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown>
  extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({
    securityWorker,
    secure,
    format,
    ...axiosConfig
  }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({
      ...axiosConfig,
      baseURL: axiosConfig.baseURL || "//localhost:8080",
    });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected mergeRequestParams(
    params1: AxiosRequestConfig,
    params2?: AxiosRequestConfig,
  ): AxiosRequestConfig {
    const method = params1.method || (params2 && params2.method);

    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...((method &&
          this.instance.defaults.headers[
            method.toLowerCase() as keyof HeadersDefaults
          ]) ||
          {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected stringifyFormItem(formItem: unknown) {
    if (typeof formItem === "object" && formItem !== null) {
      return JSON.stringify(formItem);
    } else {
      return `${formItem}`;
    }
  }

  protected createFormData(input: Record<string, unknown>): FormData {
    if (input instanceof FormData) {
      return input;
    }
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      const propertyContent: any[] =
        property instanceof Array ? property : [property];

      for (const formItem of propertyContent) {
        const isFileType = formItem instanceof Blob || formItem instanceof File;
        formData.append(
          key,
          isFileType ? formItem : this.stringifyFormItem(formItem),
        );
      }

      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = format || this.format || undefined;

    if (
      type === ContentType.FormData &&
      body &&
      body !== null &&
      typeof body === "object"
    ) {
      body = this.createFormData(body as Record<string, unknown>);
    }

    if (
      type === ContentType.Text &&
      body &&
      body !== null &&
      typeof body !== "string"
    ) {
      body = JSON.stringify(body);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type ? { "Content-Type": type } : {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title Kyooar API
 * @version 1.0
 * @termsOfService http://swagger.io/terms/
 * @baseUrl //localhost:8080
 * @contact API Support <justdevelopitmx@proton.me>
 *
 * Organization feedback management system API
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  api = {
    /**
     * @description Generate AI questions and create a complete questionnaire for a product
     *
     * @tags questionnaires, ai
     * @name V1AiGenerateQuestionnaireCreate
     * @summary Generate and save AI questionnaire
     * @request POST:/api/v1/ai/generate-questionnaire/{productId}
     * @secure
     */
    v1AiGenerateQuestionnaireCreate: (
      productId: string,
      questionnaire: ModelsGenerateQuestionnaireRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: ModelsQuestionnaire;
        },
        ResponseResponse
      >({
        path: `/api/v1/ai/generate-questionnaire/${productId}`,
        method: "POST",
        body: questionnaire,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Generate AI-powered questions for a specific product
     *
     * @tags questionnaires, ai
     * @name V1AiGenerateQuestionsCreate
     * @summary Generate AI questions
     * @request POST:/api/v1/ai/generate-questions/{productId}
     * @secure
     */
    v1AiGenerateQuestionsCreate: (
      productId: string,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: ModelsGeneratedQuestion[];
        },
        ResponseResponse
      >({
        path: `/api/v1/ai/generate-questions/${productId}`,
        method: "POST",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get basic analytics metrics for the dashboard including satisfaction, recommendation rate, and recent feedback
     *
     * @tags analytics
     * @name V1AnalyticsDashboardDetail
     * @summary Get dashboard metrics
     * @request GET:/api/v1/analytics/dashboard/{organizationId}
     * @secure
     */
    v1AnalyticsDashboardDetail: (
      organizationId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/analytics/dashboard/${organizationId}`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get comprehensive analytics data for a organization including ratings, feedback counts, and product performance
     *
     * @tags analytics
     * @name V1AnalyticsOrganizationsDetail
     * @summary Get organization analytics
     * @request GET:/api/v1/analytics/organizations/{organizationId}
     * @secure
     */
    v1AnalyticsOrganizationsDetail: (
      organizationId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/analytics/organizations/${organizationId}`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get pre-aggregated chart data for all questions in a organization with optional filters
     *
     * @tags analytics
     * @name V1AnalyticsOrganizationsChartsList
     * @summary Get organization chart data
     * @request GET:/api/v1/analytics/organizations/{organizationId}/charts
     * @secure
     */
    v1AnalyticsOrganizationsChartsList: (
      organizationId: string,
      query?: {
        /** Start date (YYYY-MM-DD) */
        date_from?: string;
        /** End date (YYYY-MM-DD) */
        date_to?: string;
        /** Filter by specific product ID */
        product_id?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/analytics/organizations/${organizationId}/charts`,
        method: "GET",
        query: query,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get detailed analytics data for a specific product including ratings, feedback count, and recent feedback
     *
     * @tags analytics
     * @name V1AnalyticsProductsDetail
     * @summary Get product analytics
     * @request GET:/api/v1/analytics/products/{productId}
     * @secure
     */
    v1AnalyticsProductsDetail: (
      productId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/analytics/products/${productId}`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get detailed insights for a specific product including question-level analytics
     *
     * @tags analytics
     * @name V1AnalyticsProductsInsightsList
     * @summary Get product insights
     * @request GET:/api/v1/analytics/products/{productId}/insights
     * @secure
     */
    v1AnalyticsProductsInsightsList: (
      productId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/analytics/products/${productId}/insights`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Cancel a pending account deactivation request
     *
     * @tags auth
     * @name V1AuthCancelDeactivationCreate
     * @summary Cancel account deactivation
     * @request POST:/api/v1/auth/cancel-deactivation
     * @secure
     */
    v1AuthCancelDeactivationCreate: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/cancel-deactivation`,
        method: "POST",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Request to change the account email address
     *
     * @tags auth
     * @name V1AuthChangeEmailCreate
     * @summary Request email change
     * @request POST:/api/v1/auth/change-email
     * @secure
     */
    v1AuthChangeEmailCreate: (
      request: HandlersChangeEmailRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/change-email`,
        method: "POST",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Confirm email change using the token sent to the new email
     *
     * @tags auth
     * @name V1AuthConfirmEmailChangeCreate
     * @summary Confirm email change
     * @request POST:/api/v1/auth/confirm-email-change
     */
    v1AuthConfirmEmailChangeCreate: (
      request: HandlersConfirmEmailChangeRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/confirm-email-change`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Request to deactivate the account with a 15-day grace period
     *
     * @tags auth
     * @name V1AuthDeactivateCreate
     * @summary Request account deactivation
     * @request POST:/api/v1/auth/deactivate
     * @secure
     */
    v1AuthDeactivateCreate: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, any>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/deactivate`,
        method: "POST",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Send password reset email to the specified email address
     *
     * @tags auth
     * @name V1AuthForgotPasswordCreate
     * @summary Send password reset email
     * @request POST:/api/v1/auth/forgot-password
     */
    v1AuthForgotPasswordCreate: (
      request: HandlersPasswordResetRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/forgot-password`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Authenticate and get JWT token
     *
     * @tags auth
     * @name V1AuthLoginCreate
     * @summary Login to account
     * @request POST:/api/v1/auth/login
     */
    v1AuthLoginCreate: (
      request: HandlersLoginRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: HandlersAuthResponse;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/login`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update user profile information including company name and personal details
     *
     * @tags auth
     * @name V1AuthProfileUpdate
     * @summary Update user profile
     * @request PUT:/api/v1/auth/profile
     * @secure
     */
    v1AuthProfileUpdate: (
      request: HandlersUpdateProfileRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: any;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/profile`,
        method: "PUT",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Refresh an existing JWT token to get a new one
     *
     * @tags auth
     * @name V1AuthRefreshCreate
     * @summary Refresh JWT token
     * @request POST:/api/v1/auth/refresh
     * @secure
     */
    v1AuthRefreshCreate: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/refresh`,
        method: "POST",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new organization owner account
     *
     * @tags auth
     * @name V1AuthRegisterCreate
     * @summary Register a new account
     * @request POST:/api/v1/auth/register
     */
    v1AuthRegisterCreate: (
      request: HandlersRegisterRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: any;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/register`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Resend verification email to the specified email address (public endpoint)
     *
     * @tags auth
     * @name V1AuthResendVerificationCreate
     * @summary Resend email verification
     * @request POST:/api/v1/auth/resend-verification
     */
    v1AuthResendVerificationCreate: (
      request: HandlersResendVerificationRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/resend-verification`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Reset password using reset token
     *
     * @tags auth
     * @name V1AuthResetPasswordCreate
     * @summary Reset password
     * @request POST:/api/v1/auth/reset-password
     */
    v1AuthResetPasswordCreate: (
      request: HandlersResetPasswordRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/reset-password`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Send verification email to the authenticated account
     *
     * @tags auth
     * @name V1AuthSendVerificationCreate
     * @summary Send email verification
     * @request POST:/api/v1/auth/send-verification
     * @secure
     */
    v1AuthSendVerificationCreate: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/send-verification`,
        method: "POST",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Verify email address using verification token
     *
     * @tags auth
     * @name V1AuthVerifyEmailList
     * @summary Verify email address
     * @request GET:/api/v1/auth/verify-email
     */
    v1AuthVerifyEmailList: (
      query: {
        /** Verification token */
        token: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/verify-email`,
        method: "GET",
        query: query,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all organizations for the authenticated account
     *
     * @tags organizations
     * @name V1OrganizationsList
     * @summary Get all organizations
     * @request GET:/api/v1/organizations
     * @secure
     */
    v1OrganizationsList: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: ModelsOrganization[];
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new organization for the authenticated account
     *
     * @tags organizations
     * @name V1OrganizationsCreate
     * @summary Create a new organization
     * @request POST:/api/v1/organizations
     * @secure
     */
    v1OrganizationsCreate: (
      request: HandlersCreateOrganizationRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: ModelsOrganization;
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations`,
        method: "POST",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get a specific organization by its ID
     *
     * @tags organizations
     * @name V1OrganizationsDetail
     * @summary Get organization by ID
     * @request GET:/api/v1/organizations/{id}
     * @secure
     */
    v1OrganizationsDetail: (id: string, params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: ModelsOrganization;
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${id}`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update a organization's information
     *
     * @tags organizations
     * @name V1OrganizationsUpdate
     * @summary Update organization
     * @request PUT:/api/v1/organizations/{id}
     * @secure
     */
    v1OrganizationsUpdate: (
      id: string,
      updates: Record<string, any>,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${id}`,
        method: "PUT",
        body: updates,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete a organization from the system
     *
     * @tags organizations
     * @name V1OrganizationsDelete
     * @summary Delete organization
     * @request DELETE:/api/v1/organizations/{id}
     * @secure
     */
    v1OrganizationsDelete: (id: string, params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${id}`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get feedback analytics and statistics for a organization
     *
     * @tags feedback
     * @name V1OrganizationsAnalyticsList
     * @summary Get feedback statistics
     * @request GET:/api/v1/organizations/{organizationId}/analytics
     * @secure
     */
    v1OrganizationsAnalyticsList: (
      organizationId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/organizations/${organizationId}/analytics`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all feedback for a specific organization with pagination and optional filters
     *
     * @tags feedback
     * @name V1OrganizationsFeedbackList
     * @summary Get organization feedback with filters
     * @request GET:/api/v1/organizations/{organizationId}/feedback
     * @secure
     */
    v1OrganizationsFeedbackList: (
      organizationId: string,
      query?: {
        /** Page number (default: 1) */
        page?: number;
        /** Items per page (default: 20, max: 100) */
        limit?: number;
        /** Search in comments, customer name, or email */
        search?: string;
        /** Minimum rating (1-5) */
        rating_min?: number;
        /** Maximum rating (1-5) */
        rating_max?: number;
        /** Start date (YYYY-MM-DD format) */
        date_from?: string;
        /** End date (YYYY-MM-DD format) */
        date_to?: string;
        /** Filter by specific product ID */
        product_id?: string;
        /** Filter by completion status */
        is_complete?: boolean;
      },
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/organizations/${organizationId}/feedback`,
        method: "GET",
        query: query,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all products for a specific organization
     *
     * @tags products
     * @name V1OrganizationsProductsList
     * @summary Get products by organization
     * @request GET:/api/v1/organizations/{organizationId}/products
     * @secure
     */
    v1OrganizationsProductsList: (
      organizationId: string,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: KyooarInternalMenuModelsProduct[];
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${organizationId}/products`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new product for a organization
     *
     * @tags products
     * @name V1OrganizationsProductsCreate
     * @summary Create a new product
     * @request POST:/api/v1/organizations/{organizationId}/products
     * @secure
     */
    v1OrganizationsProductsCreate: (
      organizationId: string,
      product: HandlersCreateProductRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: KyooarInternalMenuModelsProduct;
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${organizationId}/products`,
        method: "POST",
        body: product,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all feedback questions for a specific product
     *
     * @tags questions
     * @name V1OrganizationsProductsQuestionsList
     * @summary Get questions for a product
     * @request GET:/api/v1/organizations/{organizationId}/products/{productId}/questions
     * @secure
     */
    v1OrganizationsProductsQuestionsList: (
      organizationId: string,
      productId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/api/v1/organizations/${organizationId}/products/${productId}/questions`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Add a new feedback question to a specific product
     *
     * @tags questions
     * @name V1OrganizationsProductsQuestionsCreate
     * @summary Add a question to a product
     * @request POST:/api/v1/organizations/{organizationId}/products/{productId}/questions
     * @secure
     */
    v1OrganizationsProductsQuestionsCreate: (
      organizationId: string,
      productId: string,
      question: ModelsCreateQuestionRequest,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/api/v1/organizations/${organizationId}/products/${productId}/questions`,
        method: "POST",
        body: question,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Reorder questions for a specific product
     *
     * @tags questions
     * @name V1OrganizationsProductsQuestionsReorderCreate
     * @summary Reorder questions
     * @request POST:/api/v1/organizations/{organizationId}/products/{productId}/questions/reorder
     * @secure
     */
    v1OrganizationsProductsQuestionsReorderCreate: (
      organizationId: string,
      productId: string,
      order: string[],
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/api/v1/organizations/${organizationId}/products/${productId}/questions/reorder`,
        method: "POST",
        body: order,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get details of a specific question
     *
     * @tags questions
     * @name V1OrganizationsProductsQuestionsDetail
     * @summary Get a specific question
     * @request GET:/api/v1/organizations/{organizationId}/products/{productId}/questions/{questionId}
     * @secure
     */
    v1OrganizationsProductsQuestionsDetail: (
      organizationId: string,
      productId: string,
      questionId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/api/v1/organizations/${organizationId}/products/${productId}/questions/${questionId}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Update an existing question for a product
     *
     * @tags questions
     * @name V1OrganizationsProductsQuestionsUpdate
     * @summary Update a question
     * @request PUT:/api/v1/organizations/{organizationId}/products/{productId}/questions/{questionId}
     * @secure
     */
    v1OrganizationsProductsQuestionsUpdate: (
      organizationId: string,
      productId: string,
      questionId: string,
      question: ModelsUpdateQuestionRequest,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/api/v1/organizations/${organizationId}/products/${productId}/questions/${questionId}`,
        method: "PUT",
        body: question,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete a feedback question from a product
     *
     * @tags questions
     * @name V1OrganizationsProductsQuestionsDelete
     * @summary Delete a question
     * @request DELETE:/api/v1/organizations/{organizationId}/products/{productId}/questions/{questionId}
     * @secure
     */
    v1OrganizationsProductsQuestionsDelete: (
      organizationId: string,
      productId: string,
      questionId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/api/v1/organizations/${organizationId}/products/${productId}/questions/${questionId}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all QR codes for a specific organization
     *
     * @tags qr-codes
     * @name V1OrganizationsQrCodesList
     * @summary Get QR codes by organization
     * @request GET:/api/v1/organizations/{organizationId}/qr-codes
     * @secure
     */
    v1OrganizationsQrCodesList: (
      organizationId: string,
      params: RequestParams = {},
    ) =>
      this.request<HandlersQRCodeListResponse, ResponseResponse>({
        path: `/api/v1/organizations/${organizationId}/qr-codes`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Generate a new QR code for a organization
     *
     * @tags qr-codes
     * @name V1OrganizationsQrCodesCreate
     * @summary Generate QR code
     * @request POST:/api/v1/organizations/{organizationId}/qr-codes
     * @secure
     */
    v1OrganizationsQrCodesCreate: (
      organizationId: string,
      qr_code: HandlersGenerateQRCodeRequest,
      params: RequestParams = {},
    ) =>
      this.request<HandlersGenerateQRCodeResponse, ResponseResponse>({
        path: `/api/v1/organizations/${organizationId}/qr-codes`,
        method: "POST",
        body: qr_code,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all questionnaires for a organization
     *
     * @tags questionnaires
     * @name V1OrganizationsQuestionnairesList
     * @summary List questionnaires
     * @request GET:/api/v1/organizations/{organizationId}/questionnaires
     * @secure
     */
    v1OrganizationsQuestionnairesList: (
      organizationId: string,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: ModelsQuestionnaire[];
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${organizationId}/questionnaires`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new questionnaire for a organization
     *
     * @tags questionnaires
     * @name V1OrganizationsQuestionnairesCreate
     * @summary Create questionnaire
     * @request POST:/api/v1/organizations/{organizationId}/questionnaires
     * @secure
     */
    v1OrganizationsQuestionnairesCreate: (
      organizationId: string,
      questionnaire: ModelsCreateQuestionnaireRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: ModelsQuestionnaire;
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${organizationId}/questionnaires`,
        method: "POST",
        body: questionnaire,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get a specific questionnaire by ID
     *
     * @tags questionnaires
     * @name V1OrganizationsQuestionnairesDetail
     * @summary Get questionnaire
     * @request GET:/api/v1/organizations/{organizationId}/questionnaires/{id}
     * @secure
     */
    v1OrganizationsQuestionnairesDetail: (
      organizationId: string,
      id: string,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: ModelsQuestionnaire;
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${organizationId}/questionnaires/${id}`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update an existing questionnaire
     *
     * @tags questionnaires
     * @name V1OrganizationsQuestionnairesUpdate
     * @summary Update questionnaire
     * @request PUT:/api/v1/organizations/{organizationId}/questionnaires/{id}
     * @secure
     */
    v1OrganizationsQuestionnairesUpdate: (
      organizationId: string,
      id: string,
      questionnaire: ModelsQuestionnaire,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: ModelsQuestionnaire;
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${organizationId}/questionnaires/${id}`,
        method: "PUT",
        body: questionnaire,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete a questionnaire
     *
     * @tags questionnaires
     * @name V1OrganizationsQuestionnairesDelete
     * @summary Delete questionnaire
     * @request DELETE:/api/v1/organizations/{organizationId}/questionnaires/{id}
     * @secure
     */
    v1OrganizationsQuestionnairesDelete: (
      organizationId: string,
      id: string,
      params: RequestParams = {},
    ) =>
      this.request<ResponseResponse, ResponseResponse>({
        path: `/api/v1/organizations/${organizationId}/questionnaires/${id}`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get list of product IDs that have questions for a organization
     *
     * @tags questions
     * @name V1OrganizationsQuestionsProductsWithQuestionsList
     * @summary Get products that have questions
     * @request GET:/api/v1/organizations/{organizationId}/questions/products-with-questions
     * @secure
     */
    v1OrganizationsQuestionsProductsWithQuestionsList: (
      organizationId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/api/v1/organizations/${organizationId}/questions/products-with-questions`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a payment checkout session for a subscription plan
     *
     * @tags payment
     * @name V1PaymentCheckoutCreate
     * @summary Create a checkout session
     * @request POST:/api/v1/payment/checkout
     * @secure
     */
    v1PaymentCheckoutCreate: (
      request: HandlersCreateCheckoutRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: HandlersCheckoutResponse;
        },
        ResponseResponse
      >({
        path: `/api/v1/payment/checkout`,
        method: "POST",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Complete a checkout session after payment
     *
     * @tags payment
     * @name V1PaymentCheckoutCompleteCreate
     * @summary Complete a checkout session
     * @request POST:/api/v1/payment/checkout/complete
     */
    v1PaymentCheckoutCompleteCreate: (
      request: HandlersCompleteCheckoutRequest,
      params: RequestParams = {},
    ) =>
      this.request<ResponseResponse, ResponseResponse>({
        path: `/api/v1/payment/checkout/complete`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get user's invoice history
     *
     * @tags payment
     * @name V1PaymentInvoicesList
     * @summary Get invoices
     * @request GET:/api/v1/payment/invoices
     * @secure
     */
    v1PaymentInvoicesList: (
      query?: {
        /**
         * Limit number of invoices
         * @default 10
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: HandlersInvoiceResponse[];
        },
        ResponseResponse
      >({
        path: `/api/v1/payment/invoices`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Get list of user's payment methods
     *
     * @tags payment
     * @name V1PaymentMethodsList
     * @summary List payment methods
     * @request GET:/api/v1/payment/methods
     * @secure
     */
    v1PaymentMethodsList: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: HandlersPaymentMethodResponse[];
        },
        ResponseResponse
      >({
        path: `/api/v1/payment/methods`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Set a payment method as default
     *
     * @tags payment
     * @name V1PaymentMethodsDefaultCreate
     * @summary Set default payment method
     * @request POST:/api/v1/payment/methods/default
     * @secure
     */
    v1PaymentMethodsDefaultCreate: (
      request: HandlersSetDefaultPaymentRequest,
      params: RequestParams = {},
    ) =>
      this.request<ResponseResponse, ResponseResponse>({
        path: `/api/v1/payment/methods/default`,
        method: "POST",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a customer portal session for self-service subscription management
     *
     * @tags payment
     * @name V1PaymentPortalCreate
     * @summary Create customer portal session
     * @request POST:/api/v1/payment/portal
     * @secure
     */
    v1PaymentPortalCreate: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: HandlersPortalResponse;
        },
        ResponseResponse
      >({
        path: `/api/v1/payment/portal`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Handle webhook events from payment provider
     *
     * @tags payment
     * @name V1PaymentWebhookCreate
     * @summary Handle payment webhook
     * @request POST:/api/v1/payment/webhook
     */
    v1PaymentWebhookCreate: (params: RequestParams = {}) =>
      this.request<ResponseResponse, ResponseResponse>({
        path: `/api/v1/payment/webhook`,
        method: "POST",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Retrieve all available subscription plans with their features and pricing
     *
     * @tags subscription
     * @name V1PlansList
     * @summary Get available subscription plans
     * @request GET:/api/v1/plans
     */
    v1PlansList: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: ModelsSubscriptionPlan[];
        },
        ResponseResponse
      >({
        path: `/api/v1/plans`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get a specific product by its ID
     *
     * @tags products
     * @name V1ProductsDetail
     * @summary Get product by ID
     * @request GET:/api/v1/products/{id}
     * @secure
     */
    v1ProductsDetail: (id: string, params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: KyooarInternalMenuModelsProduct;
        },
        ResponseResponse
      >({
        path: `/api/v1/products/${id}`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update a product's information
     *
     * @tags products
     * @name V1ProductsUpdate
     * @summary Update a product
     * @request PUT:/api/v1/products/{id}
     * @secure
     */
    v1ProductsUpdate: (
      id: string,
      updates: Record<string, any>,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/products/${id}`,
        method: "PUT",
        body: updates,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete a product from the system
     *
     * @tags products
     * @name V1ProductsDelete
     * @summary Delete a product
     * @request DELETE:/api/v1/products/{id}
     * @secure
     */
    v1ProductsDelete: (id: string, params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/products/${id}`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Submit customer feedback for a product
     *
     * @tags public
     * @name V1PublicFeedbackCreate
     * @summary Submit feedback
     * @request POST:/api/v1/public/feedback
     */
    v1PublicFeedbackCreate: (
      feedback: ModelsFeedback,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/public/feedback`,
        method: "POST",
        body: feedback,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get public menu for a organization
     *
     * @tags public
     * @name V1PublicOrganizationMenuList
     * @summary Get organization menu
     * @request GET:/api/v1/public/organization/{id}/menu
     */
    v1PublicOrganizationMenuList: (id: string, params: RequestParams = {}) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/public/organization/${id}/menu`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all feedback questions for a specific product (public access for customer feedback)
     *
     * @tags public
     * @name V1PublicOrganizationProductsQuestionsList
     * @summary Get questions for a product
     * @request GET:/api/v1/public/organization/{organizationId}/products/{productId}/questions
     */
    v1PublicOrganizationProductsQuestionsList: (
      organizationId: string,
      productId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/public/organization/${organizationId}/products/${productId}/questions`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description Get all products that have feedback questions for a organization (public access for QR code scans)
     *
     * @tags public
     * @name V1PublicOrganizationQuestionsProductsWithQuestionsList
     * @summary Get products with questions
     * @request GET:/api/v1/public/organization/{organizationId}/questions/products-with-questions
     */
    v1PublicOrganizationQuestionsProductsWithQuestionsList: (
      organizationId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/public/organization/${organizationId}/questions/products-with-questions`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description Validate a QR code and return associated data
     *
     * @tags public
     * @name V1PublicQrDetail
     * @summary Validate QR code
     * @request GET:/api/v1/public/qr/{code}
     */
    v1PublicQrDetail: (code: string, params: RequestParams = {}) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/public/qr/${code}`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get questionnaire for a specific product
     *
     * @tags public
     * @name V1PublicQuestionnaireDetail
     * @summary Get questionnaire
     * @request GET:/api/v1/public/questionnaire/{organizationId}/{productId}
     */
    v1PublicQuestionnaireDetail: (
      organizationId: string,
      productId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/public/questionnaire/${organizationId}/${productId}`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete a QR code from the system
     *
     * @tags qr-codes
     * @name V1QrCodesDelete
     * @summary Delete QR code
     * @request DELETE:/api/v1/qr-codes/{id}
     * @secure
     */
    v1QrCodesDelete: (id: string, params: RequestParams = {}) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/qr-codes/${id}`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update QR code details like active status, label, or location
     *
     * @tags qr-codes
     * @name V1QrCodesPartialUpdate
     * @summary Update QR code
     * @request PATCH:/api/v1/qr-codes/{id}
     * @secure
     */
    v1QrCodesPartialUpdate: (
      id: string,
      qr_code: HandlersUpdateQRCodeRequest,
      params: RequestParams = {},
    ) =>
      this.request<HandlersUpdateQRCodeResponse, ResponseResponse>({
        path: `/api/v1/qr-codes/${id}`,
        method: "PATCH",
        body: qr_code,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Accept a team invitation using the invitation token
     *
     * @tags team
     * @name V1TeamAcceptInviteCreate
     * @summary Accept team invitation
     * @request POST:/api/v1/team/accept-invite
     */
    v1TeamAcceptInviteCreate: (
      request: HandlersAcceptInviteRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/team/accept-invite`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all team members for the account
     *
     * @tags team
     * @name V1TeamMembersList
     * @summary List team members
     * @request GET:/api/v1/team/members
     * @secure
     */
    v1TeamMembersList: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: ModelsTeamMember[];
        },
        ResponseResponse
      >({
        path: `/api/v1/team/members`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Invite a new team member to the account
     *
     * @tags team
     * @name V1TeamMembersInviteCreate
     * @summary Invite team member
     * @request POST:/api/v1/team/members/invite
     * @secure
     */
    v1TeamMembersInviteCreate: (
      request: HandlersInviteMemberRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: ModelsTeamMember;
        },
        ResponseResponse
      >({
        path: `/api/v1/team/members/invite`,
        method: "POST",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Remove a team member from the account
     *
     * @tags team
     * @name V1TeamMembersDelete
     * @summary Remove team member
     * @request DELETE:/api/v1/team/members/{id}
     * @secure
     */
    v1TeamMembersDelete: (id: string, params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/team/members/${id}`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Resend an invitation email to a pending team member
     *
     * @tags team
     * @name V1TeamMembersResendInvitationCreate
     * @summary Resend team invitation
     * @request POST:/api/v1/team/members/{id}/resend-invitation
     * @secure
     */
    v1TeamMembersResendInvitationCreate: (
      id: string,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/team/members/${id}/resend-invitation`,
        method: "POST",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update the role of a team member
     *
     * @tags team
     * @name V1TeamMembersRoleUpdate
     * @summary Update team member role
     * @request PUT:/api/v1/team/members/{id}/role
     * @secure
     */
    v1TeamMembersRoleUpdate: (
      id: string,
      request: HandlersUpdateRoleRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/team/members/${id}/role`,
        method: "PUT",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Check if the authenticated user can create more organizations based on their subscription plan
     *
     * @tags subscription
     * @name V1UserCanCreateOrganizationList
     * @summary Check if user can create more organizations
     * @request GET:/api/v1/user/can-create-organization
     * @secure
     */
    v1UserCanCreateOrganizationList: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: ServicesPermissionResponse;
        },
        ResponseResponse
      >({
        path: `/api/v1/user/can-create-organization`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Retrieve the current subscription details for the authenticated user
     *
     * @tags subscription
     * @name V1UserSubscriptionList
     * @summary Get user's current subscription
     * @request GET:/api/v1/user/subscription
     * @secure
     */
    v1UserSubscriptionList: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: ModelsSubscription;
        },
        ResponseResponse
      >({
        path: `/api/v1/user/subscription`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new subscription for the authenticated user
     *
     * @tags subscription
     * @name V1UserSubscriptionCreate
     * @summary Create a new subscription
     * @request POST:/api/v1/user/subscription
     * @secure
     */
    v1UserSubscriptionCreate: (
      request: HandlersCreateSubscriptionRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: ModelsSubscription;
        },
        ResponseResponse
      >({
        path: `/api/v1/user/subscription`,
        method: "POST",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Cancel the current subscription for the authenticated user
     *
     * @tags subscription
     * @name V1UserSubscriptionDelete
     * @summary Cancel user's subscription
     * @request DELETE:/api/v1/user/subscription
     * @secure
     */
    v1UserSubscriptionDelete: (params: RequestParams = {}) =>
      this.request<ResponseResponse, ResponseResponse>({
        path: `/api/v1/user/subscription`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Retrieve the current subscription usage details for the authenticated user
     *
     * @tags subscription
     * @name V1UserSubscriptionUsageList
     * @summary Get user's current subscription usage
     * @request GET:/api/v1/user/subscription/usage
     * @secure
     */
    v1UserSubscriptionUsageList: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: ModelsSubscriptionUsage;
        },
        ResponseResponse
      >({
        path: `/api/v1/user/subscription/usage`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  organizations = {
    /**
     * @description Add a new question to an existing questionnaire
     *
     * @tags questionnaires
     * @name QuestionnairesQuestionsCreate
     * @summary Add a question to a questionnaire
     * @request POST:/organizations/{organizationId}/questionnaires/{id}/questions
     * @secure
     */
    questionnairesQuestionsCreate: (
      id: string,
      organizationId: string,
      question: ModelsQuestion,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/organizations/${organizationId}/questionnaires/${id}/questions`,
        method: "POST",
        body: question,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update an existing question in a questionnaire
     *
     * @tags questionnaires
     * @name QuestionnairesQuestionsUpdate
     * @summary Update a question
     * @request PUT:/organizations/{organizationId}/questionnaires/{id}/questions/{questionId}
     * @secure
     */
    questionnairesQuestionsUpdate: (
      id: string,
      questionId: string,
      organizationId: string,
      question: ModelsQuestion,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/organizations/${organizationId}/questionnaires/${id}/questions/${questionId}`,
        method: "PUT",
        body: question,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete a question from a questionnaire
     *
     * @tags questionnaires
     * @name QuestionnairesQuestionsDelete
     * @summary Delete a question
     * @request DELETE:/organizations/{organizationId}/questionnaires/{id}/questions/{questionId}
     * @secure
     */
    questionnairesQuestionsDelete: (
      id: string,
      questionId: string,
      organizationId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/organizations/${organizationId}/questionnaires/${id}/questions/${questionId}`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Reorder questions in a questionnaire
     *
     * @tags questionnaires
     * @name QuestionnairesReorderCreate
     * @summary Reorder questions
     * @request POST:/organizations/{organizationId}/questionnaires/{id}/reorder
     * @secure
     */
    questionnairesReorderCreate: (
      id: string,
      organizationId: string,
      order: string[],
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/organizations/${organizationId}/questionnaires/${id}/reorder`,
        method: "POST",
        body: order,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
}
