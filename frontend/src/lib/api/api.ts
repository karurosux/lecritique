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

export enum SubscriptionmodelSubscriptionStatus {
  SubscriptionActive = "active",
  SubscriptionPending = "pending",
  SubscriptionCanceled = "canceled",
  SubscriptionExpired = "expired",
}

export enum QrcodemodelQRCodeType {
  QRCodeTypeTable = "table",
  QRCodeTypeLocation = "location",
  QRCodeTypeTakeaway = "takeaway",
  QRCodeTypeDelivery = "delivery",
  QRCodeTypeGeneral = "general",
}

export enum ModelsMemberRole {
  RoleOwner = "OWNER",
  RoleAdmin = "ADMIN",
  RoleManager = "MANAGER",
  RoleViewer = "VIEWER",
}

export enum FeedbackmodelQuestionType {
  QuestionTypeRating = "rating",
  QuestionTypeScale = "scale",
  QuestionTypeMultiChoice = "multi_choice",
  QuestionTypeSingleChoice = "single_choice",
  QuestionTypeText = "text",
  QuestionTypeYesNo = "yes_no",
}

export interface AnalyticsmodelChoiceInfo {
  choice?: string;
  count?: number;
}

export interface AnalyticsmodelChoiceSeriesData {
  choice?: string;
  points?: AnalyticsmodelTimeSeriesPoint[];
  statistics?: AnalyticsmodelTimeSeriesStats;
}

export interface AnalyticsmodelComparisonInsight {
  change?: number;
  message?: string;
  metric_type?: string;
  recommendation?: string;
  severity?: string;
  type?: string;
}

export interface AnalyticsmodelComparisonRequest {
  /** @minItems 1 */
  metric_types: string[];
  organization_id: string;
  period1_end: string;
  period1_start: string;
  period2_end: string;
  period2_start: string;
  product_id?: string;
  question_id?: string;
}

export interface AnalyticsmodelComparisonResponse {
  comparisons?: AnalyticsmodelTimeSeriesComparison[];
  insights?: AnalyticsmodelComparisonInsight[];
  request?: AnalyticsmodelComparisonRequest;
}

export interface AnalyticsmodelDateRange {
  end?: string;
  start?: string;
}

export interface AnalyticsmodelTimePeriodMetrics {
  average?: number;
  choice_distribution?: Record<string, number>;
  count?: number;
  data_points?: AnalyticsmodelTimeSeriesPoint[];
  end_date?: string;
  max?: number;
  min?: number;
  most_popular_choice?: AnalyticsmodelChoiceInfo;
  start_date?: string;
  top_choices?: AnalyticsmodelChoiceInfo[];
  value?: number;
}

export interface AnalyticsmodelTimeSeriesComparison {
  change?: number;
  change_percent?: number;
  metadata?: string;
  metric_name?: string;
  metric_type?: string;
  period1?: AnalyticsmodelTimePeriodMetrics;
  period2?: AnalyticsmodelTimePeriodMetrics;
  trend?: string;
}

export interface AnalyticsmodelTimeSeriesData {
  choice_series?: AnalyticsmodelChoiceSeriesData[];
  metadata?: Record<string, any>;
  metric_name?: string;
  metric_type?: string;
  points?: AnalyticsmodelTimeSeriesPoint[];
  product_id?: string;
  product_name?: string;
  statistics?: AnalyticsmodelTimeSeriesStats;
}

export interface AnalyticsmodelTimeSeriesPoint {
  count?: number;
  timestamp?: string;
  value?: number;
}

export interface AnalyticsmodelTimeSeriesRequest {
  end_date: string;
  granularity: "hourly" | "daily" | "weekly" | "monthly";
  /** @minItems 1 */
  metric_types: string[];
  organization_id: string;
  product_id?: string;
  question_id?: string;
  start_date: string;
}

export interface AnalyticsmodelTimeSeriesResponse {
  request?: AnalyticsmodelTimeSeriesRequest;
  series?: AnalyticsmodelTimeSeriesData[];
  summary?: AnalyticsmodelTimeSeriesSummary;
}

export interface AnalyticsmodelTimeSeriesStats {
  average?: number;
  count?: number;
  max?: number;
  min?: number;
  total?: number;
  trend_direction?: string;
  trend_strength?: number;
}

export interface AnalyticsmodelTimeSeriesSummary {
  date_range?: AnalyticsmodelDateRange;
  granularity?: string;
  metrics_summary?: Record<string, any>;
  total_data_points?: number;
}

export interface AuthmodelAcceptInvitationRequest {
  token: string;
}

export interface AuthmodelChangeEmailRequest {
  new_email: string;
}

export interface AuthmodelConfirmEmailChangeRequest {
  token: string;
}

export interface AuthmodelDeactivationResponse {
  deactivation_date?: string;
  message?: string;
}

export interface AuthmodelInvitationResponse {
  invitation?: any;
  message?: string;
}

export interface AuthmodelInviteMemberRequest {
  email: string;
  role: "OWNER" | "ADMIN" | "MANAGER" | "VIEWER";
}

export interface AuthmodelLoginRequest {
  email: string;
  password: string;
}

export interface AuthmodelMemberListResponse {
  members?: ModelsTeamMember[];
}

export interface AuthmodelPasswordResetRequest {
  email: string;
}

export interface AuthmodelRegisterRequest {
  email: string;
  first_name?: string;
  invitation_token?: string;
  last_name?: string;
  name?: string;
  /** @minLength 8 */
  password: string;
}

export interface AuthmodelResendVerificationRequest {
  email: string;
}

export interface AuthmodelResetPasswordRequest {
  /** @minLength 8 */
  new_password: string;
  token: string;
}

export interface AuthmodelTokenResponse {
  token?: string;
}

export interface AuthmodelUpdateProfileRequest {
  /** @minLength 1 */
  name?: string;
  phone?: string;
}

export interface AuthmodelUpdateRoleRequest {
  role: "OWNER" | "ADMIN" | "MANAGER" | "VIEWER";
}

export interface FeedbackmodelBatchQuestionsRequest {
  product_ids: string[];
}

export interface FeedbackmodelCreateQuestionRequest {
  is_required?: boolean;
  max_label?: string;
  max_value?: number;
  min_label?: string;
  min_value?: number;
  options?: string[];
  text: string;
  type: FeedbackmodelQuestionType;
}

export interface FeedbackmodelCreateQuestionnaireRequest {
  description?: string;
  is_default?: boolean;
  name: string;
  product_id?: string;
}

export interface FeedbackmodelDeviceInfo {
  browser?: string;
  ip?: string;
  platform?: string;
  user_agent?: string;
}

export interface FeedbackmodelFeedback {
  created_at?: string;
  customer_email?: string;
  customer_name?: string;
  customer_phone?: string;
  device_info?: FeedbackmodelDeviceInfo;
  id?: string;
  is_complete?: boolean;
  organization?: OrganizationmodelOrganization;
  organization_id?: string;
  overall_rating?: number;
  product?: ModelsProduct;
  product_id?: string;
  qr_code?: QrcodemodelQRCode;
  qr_code_id?: string;
  responses?: FeedbackmodelResponse[];
  updated_at?: string;
}

export interface FeedbackmodelGenerateQuestionnaireRequest {
  description?: string;
  is_default?: boolean;
  name: string;
}

export interface FeedbackmodelGeneratedQuestion {
  max_label?: string;
  max_value?: number;
  min_label?: string;
  min_value?: number;
  options?: string[];
  text?: string;
  type?: FeedbackmodelQuestionType;
}

export interface FeedbackmodelQuestion {
  created_at?: string;
  display_order?: number;
  id?: string;
  is_required?: boolean;
  max_label?: string;
  max_value?: number;
  min_label?: string;
  min_value?: number;
  options?: string[];
  product?: ModelsProduct;
  product_id?: string;
  text?: string;
  type?: FeedbackmodelQuestionType;
  updated_at?: string;
}

export interface FeedbackmodelQuestionnaire {
  created_at?: string;
  description?: string;
  id?: string;
  is_active?: boolean;
  is_default?: boolean;
  name?: string;
  organization?: OrganizationmodelOrganization;
  organization_id?: string;
  product?: ModelsProduct;
  product_id?: string;
  questions?: FeedbackmodelQuestion[];
  updated_at?: string;
}

export interface FeedbackmodelResponse {
  answer?: any;
  question_id?: string;
  question_text?: string;
  question_type?: FeedbackmodelQuestionType;
}

export interface FeedbackmodelUpdateQuestionRequest {
  is_required?: boolean;
  max_label?: string;
  max_value?: number;
  min_label?: string;
  min_value?: number;
  options?: string[];
  text?: string;
  type?: FeedbackmodelQuestionType;
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
  subscription?: any;
  subscription_id?: string;
  team_members?: ModelsTeamMember[];
  updated_at?: string;
}

export interface ModelsProduct {
  category?: string;
  created_at?: string;
  currency?: string;
  description?: string;
  display_order?: number;
  id?: string;
  image?: string;
  is_active?: boolean;
  is_available?: boolean;
  name?: string;
  organization_id?: string;
  price?: number;
  tags?: string[];
  updated_at?: string;
}

export interface ModelsTeamMember {
  accepted_at?: string;
  account?: ModelsAccount;
  account_id?: string;
  created_at?: string;
  id?: string;
  invited_at?: string;
  invited_by?: string;
  member?: ModelsAccount;
  member_id?: string;
  role?: ModelsMemberRole;
  updated_at?: string;
}

export interface OrganizationmodelCreateOrganizationRequest {
  address?: string;
  description?: string;
  email?: string;
  name: string;
  phone?: string;
  website?: string;
}

export interface OrganizationmodelOrganization {
  account_id?: string;
  address?: string;
  created_at?: string;
  description?: string;
  email?: string;
  id?: string;
  is_active?: boolean;
  logo?: string;
  name?: string;
  phone?: string;
  settings?: OrganizationmodelSettings;
  updated_at?: string;
  website?: string;
}

export interface OrganizationmodelSettings {
  feedback_notification?: boolean;
  language?: string;
  low_rating_threshold?: number;
  timezone?: string;
}

export interface QrcodecontrollerGenerateQRCodeRequest {
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

export interface QrcodecontrollerUpdateQRCodeRequest {
  is_active?: boolean;
  /**
   * @minLength 1
   * @maxLength 100
   */
  label?: string;
  /** @maxLength 200 */
  location?: string;
}

export interface QrcodemodelQRCode {
  code?: string;
  created_at?: string;
  expires_at?: string;
  id?: string;
  is_active?: boolean;
  label?: string;
  last_scanned_at?: string;
  location?: string;
  organization?: OrganizationmodelOrganization;
  organization_id?: string;
  scans_count?: number;
  type?: QrcodemodelQRCodeType;
  updated_at?: string;
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

export interface SubscriptioncontrollerCardDetailsResponse {
  brand?: string;
  exp_month?: number;
  exp_year?: number;
  last4?: string;
}

export interface SubscriptioncontrollerCheckoutResponse {
  checkout_url?: string;
  session_id?: string;
}

export interface SubscriptioncontrollerCompleteCheckoutRequest {
  session_id: string;
}

export interface SubscriptioncontrollerCreateCheckoutRequest {
  plan_id: string;
}

export interface SubscriptioncontrollerCreateSubscriptionRequest {
  plan_id: string;
}

export interface SubscriptioncontrollerInvoiceResponse {
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

export interface SubscriptioncontrollerPaymentMethodResponse {
  card?: SubscriptioncontrollerCardDetailsResponse;
  id?: string;
  is_default?: boolean;
  type?: string;
}

export interface SubscriptioncontrollerPortalResponse {
  portal_url?: string;
}

export interface SubscriptioncontrollerSetDefaultPaymentRequest {
  payment_method_id: string;
}

export interface SubscriptioninterfacePermissionResponse {
  can_create?: boolean;
  current_count?: number;
  max_allowed?: number;
  reason?: string;
  subscription_status?: string;
}

export interface SubscriptionmodelSubscription {
  account?: ModelsAccount;
  account_id?: string;
  cancel_at?: string;
  cancelled_at?: string;
  created_at?: string;
  current_period_end?: string;
  current_period_start?: string;
  id?: string;
  plan?: SubscriptionmodelSubscriptionPlan;
  plan_id?: string;
  status?: SubscriptionmodelSubscriptionStatus;
  updated_at?: string;
}

export interface SubscriptionmodelSubscriptionPlan {
  code?: string;
  created_at?: string;
  currency?: string;
  description?: string;
  has_advanced_analytics?: boolean;
  has_basic_analytics?: boolean;
  has_custom_branding?: boolean;
  has_feedback_explorer?: boolean;
  has_priority_support?: boolean;
  id?: string;
  interval?: string;
  is_active?: boolean;
  is_visible?: boolean;
  max_feedbacks_per_month?: number;
  max_organizations?: number;
  max_qr_codes?: number;
  max_team_members?: number;
  name?: string;
  price?: number;
  trial_days?: number;
  updated_at?: string;
}

export interface SubscriptionmodelSubscriptionUsage {
  created_at?: string;
  feedbacks_count?: number;
  id?: string;
  last_updated_at?: string;
  locations_count?: number;
  organizations_count?: number;
  period_end?: string;
  period_start?: string;
  qr_codes_count?: number;
  subscription?: SubscriptionmodelSubscription;
  subscription_id?: string;
  team_members_count?: number;
  updated_at?: string;
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown>
  extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "//localhost:8080";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) =>
    fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter(
      (key) => "undefined" !== typeof query[key],
    );
    return keys
      .map((key) =>
        Array.isArray(query[key])
          ? this.addArrayQueryParam(query, key)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.JsonApi]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.Text]: (input: any) =>
      input !== null && typeof input !== "string"
        ? JSON.stringify(input)
        : input,
    [ContentType.FormData]: (input: any) => {
      if (input instanceof FormData) {
        return input;
      }

      return Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
              ? JSON.stringify(property)
              : `${property}`,
        );
        return formData;
      }, new FormData());
    },
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(
    params1: RequestParams,
    params2?: RequestParams,
  ): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (
    cancelToken: CancelToken,
  ): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(
      `${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`,
      {
        ...requestParams,
        headers: {
          ...(requestParams.headers || {}),
          ...(type && type !== ContentType.FormData
            ? { "Content-Type": type }
            : {}),
        },
        signal:
          (cancelToken
            ? this.createAbortSignal(cancelToken)
            : requestParams.signal) || null,
        body:
          typeof body === "undefined" || body === null
            ? null
            : payloadFormatter(body),
      },
    ).then(async (response) => {
      const r = response.clone() as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
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
      questionnaire: FeedbackmodelGenerateQuestionnaireRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: FeedbackmodelQuestionnaire;
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
          data?: FeedbackmodelGeneratedQuestion[];
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
     * @description Manually trigger the collection of time series metrics for analytics
     *
     * @tags analytics
     * @name V1AnalyticsOrganizationsCollectMetricsCreate
     * @summary Collect metrics for an organization
     * @request POST:/api/v1/analytics/organizations/{organizationId}/collect-metrics
     * @secure
     */
    v1AnalyticsOrganizationsCollectMetricsCreate: (
      organizationId: string,
      params: RequestParams = {},
    ) =>
      this.request<ResponseResponse, ResponseResponse>({
        path: `/api/v1/analytics/organizations/${organizationId}/collect-metrics`,
        method: "POST",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Compare metrics between two different time periods to identify trends and changes
     *
     * @tags analytics
     * @name V1AnalyticsOrganizationsCompareCreate
     * @summary Compare analytics between two time periods
     * @request POST:/api/v1/analytics/organizations/{organizationId}/compare
     * @secure
     */
    v1AnalyticsOrganizationsCompareCreate: (
      organizationId: string,
      body: AnalyticsmodelComparisonRequest,
      params: RequestParams = {},
    ) =>
      this.request<AnalyticsmodelComparisonResponse, ResponseResponse>({
        path: `/api/v1/analytics/organizations/${organizationId}/compare`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get time series data for various metrics with customizable granularity and date range
     *
     * @tags analytics
     * @name V1AnalyticsOrganizationsTimeSeriesList
     * @summary Get time series analytics data
     * @request GET:/api/v1/analytics/organizations/{organizationId}/time-series
     * @secure
     */
    v1AnalyticsOrganizationsTimeSeriesList: (
      organizationId: string,
      query: {
        /** Metric types to retrieve */
        metric_types: string[];
        /** Start date (ISO 8601) */
        start_date: string;
        /** End date (ISO 8601) */
        end_date: string;
        /** Data granularity (hourly, daily, weekly, monthly) */
        granularity: string;
        /** Filter by product ID */
        product_id?: string;
        /** Filter by question ID */
        question_id?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<AnalyticsmodelTimeSeriesResponse, ResponseResponse>({
        path: `/api/v1/analytics/organizations/${organizationId}/time-series`,
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
          data?: AuthmodelDeactivationResponse;
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
     * @description Cancel a pending account deactivation request
     *
     * @tags auth
     * @name V1AuthDeactivateCancelCreate
     * @summary Cancel account deactivation
     * @request POST:/api/v1/auth/deactivate/cancel
     * @secure
     */
    v1AuthDeactivateCancelCreate: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/deactivate/cancel`,
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
     * @name V1AuthEmailChangeCreate
     * @summary Request email change
     * @request POST:/api/v1/auth/email-change
     * @secure
     */
    v1AuthEmailChangeCreate: (
      request: AuthmodelChangeEmailRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/email-change`,
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
     * @name V1AuthEmailChangeConfirmCreate
     * @summary Confirm email change
     * @request POST:/api/v1/auth/email-change/confirm
     */
    v1AuthEmailChangeConfirmCreate: (
      request: AuthmodelConfirmEmailChangeRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/auth/email-change/confirm`,
        method: "POST",
        body: request,
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
      request: AuthmodelPasswordResetRequest,
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
      request: AuthmodelLoginRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: AuthmodelTokenResponse;
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
      request: AuthmodelUpdateProfileRequest,
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
          data?: AuthmodelTokenResponse;
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
      request: AuthmodelRegisterRequest,
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
      request: AuthmodelResendVerificationRequest,
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
      request: AuthmodelResetPasswordRequest,
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
          data?: OrganizationmodelOrganization[];
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
      request: OrganizationmodelCreateOrganizationRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: OrganizationmodelOrganization;
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
          data?: OrganizationmodelOrganization;
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
          data?: ModelsProduct[];
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
          data?: ModelsProduct;
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
     * @description Get a specific product by its ID
     *
     * @tags products
     * @name V1OrganizationsProductsDetail
     * @summary Get product by ID
     * @request GET:/api/v1/organizations/{organizationId}/products/{productId}
     * @secure
     */
    v1OrganizationsProductsDetail: (
      organizationId: string,
      productId: string,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: ModelsProduct;
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${organizationId}/products/${productId}`,
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
     * @name V1OrganizationsProductsUpdate
     * @summary Update a product
     * @request PUT:/api/v1/organizations/{organizationId}/products/{productId}
     * @secure
     */
    v1OrganizationsProductsUpdate: (
      organizationId: string,
      productId: string,
      updates: Record<string, any>,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${organizationId}/products/${productId}`,
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
     * @name V1OrganizationsProductsDelete
     * @summary Delete a product
     * @request DELETE:/api/v1/organizations/{organizationId}/products/{productId}
     * @secure
     */
    v1OrganizationsProductsDelete: (
      organizationId: string,
      productId: string,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/organizations/${organizationId}/products/${productId}`,
        method: "DELETE",
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
      question: FeedbackmodelCreateQuestionRequest,
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
      question: FeedbackmodelUpdateQuestionRequest,
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
      this.request<
        ResponseResponse & {
          data?: QrcodemodelQRCode[];
        },
        ResponseResponse
      >({
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
      qr_code: QrcodecontrollerGenerateQRCodeRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: QrcodemodelQRCode;
        },
        ResponseResponse
      >({
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
          data?: FeedbackmodelQuestionnaire[];
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
      questionnaire: FeedbackmodelCreateQuestionnaireRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: FeedbackmodelQuestionnaire;
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
          data?: FeedbackmodelQuestionnaire;
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
      questionnaire: FeedbackmodelQuestionnaire,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: FeedbackmodelQuestionnaire;
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
     * @description Get essential question fields for multiple products in a single request - returns only ID, ProductID, Text, and Type
     *
     * @tags questions
     * @name V1OrganizationsQuestionsBatchCreate
     * @summary Get questions for multiple products (optimized payload)
     * @request POST:/api/v1/organizations/{organizationId}/questions/batch
     * @secure
     */
    v1OrganizationsQuestionsBatchCreate: (
      organizationId: string,
      request: FeedbackmodelBatchQuestionsRequest,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, Record<string, any>>({
        path: `/api/v1/organizations/${organizationId}/questions/batch`,
        method: "POST",
        body: request,
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
      request: SubscriptioncontrollerCreateCheckoutRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: SubscriptioncontrollerCheckoutResponse;
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
      request: SubscriptioncontrollerCompleteCheckoutRequest,
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
          data?: SubscriptioncontrollerInvoiceResponse[];
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
          data?: SubscriptioncontrollerPaymentMethodResponse[];
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
      request: SubscriptioncontrollerSetDefaultPaymentRequest,
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
          data?: SubscriptioncontrollerPortalResponse;
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
          data?: SubscriptionmodelSubscriptionPlan[];
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
     * @description Submit customer feedback for a product
     *
     * @tags public
     * @name V1PublicFeedbackCreate
     * @summary Submit feedback
     * @request POST:/api/v1/public/feedback
     */
    v1PublicFeedbackCreate: (
      feedback: FeedbackmodelFeedback,
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
     * @description Get public products for a organization
     *
     * @tags public
     * @name V1PublicOrganizationProductsList
     * @summary Get organization products
     * @request GET:/api/v1/public/organization/{id}/products
     */
    v1PublicOrganizationProductsList: (
      id: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/public/organization/${id}/products`,
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
      this.request<
        ResponseResponse & {
          data?: QrcodemodelQRCode;
        },
        ResponseResponse
      >({
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
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
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
      qr_code: QrcodecontrollerUpdateQRCodeRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: QrcodemodelQRCode;
        },
        ResponseResponse
      >({
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
     * @name V1TeamAcceptInvitationCreate
     * @summary Accept team invitation
     * @request POST:/api/v1/team/accept-invitation
     */
    v1TeamAcceptInvitationCreate: (
      request: AuthmodelAcceptInvitationRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: any;
        },
        ResponseResponse
      >({
        path: `/api/v1/team/accept-invitation`,
        method: "POST",
        body: request,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get list of team members for the authenticated account
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
          data?: AuthmodelMemberListResponse;
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
     * @description Invite a new member to the team
     *
     * @tags team
     * @name V1TeamMembersInviteCreate
     * @summary Invite team member
     * @request POST:/api/v1/team/members/invite
     * @secure
     */
    v1TeamMembersInviteCreate: (
      request: AuthmodelInviteMemberRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: AuthmodelInvitationResponse;
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
     * @description Remove a member from the team
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
     * @description Resend invitation to a team member
     *
     * @tags team
     * @name V1TeamMembersResendInvitationCreate
     * @summary Resend invitation
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
     * @summary Update member role
     * @request PUT:/api/v1/team/members/{id}/role
     * @secure
     */
    v1TeamMembersRoleUpdate: (
      id: string,
      request: AuthmodelUpdateRoleRequest,
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
          data?: SubscriptioninterfacePermissionResponse;
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
          data?: SubscriptionmodelSubscription;
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
      request: SubscriptioncontrollerCreateSubscriptionRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: SubscriptionmodelSubscription;
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
          data?: SubscriptionmodelSubscriptionUsage;
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
      question: FeedbackmodelQuestion,
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
      question: FeedbackmodelQuestion,
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
