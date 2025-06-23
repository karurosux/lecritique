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

export type GithubComLecritiqueApiInternalMenuModelsDish = object;

export interface HandlersAuthResponse {
  account?: any;
  token?: string;
}

export interface HandlersCreateDishRequest {
  category?: string;
  currency?: string;
  description?: string;
  name: string;
  /** @min 0 */
  price?: number;
  restaurant_id: string;
}

export interface HandlersCreateRestaurantRequest {
  description?: string;
  email?: string;
  name: string;
  phone?: string;
  website?: string;
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
  restaurant_id: string;
  type: "table" | "location" | "takeaway" | "delivery" | "general";
}

export interface HandlersGenerateQRCodeResponse {
  data?: ModelsQRCode;
  success?: boolean;
}

export interface HandlersLoginRequest {
  email: string;
  password: string;
}

export interface HandlersPasswordResetRequest {
  email: string;
}

export interface HandlersQRCodeListResponse {
  data?: ModelsQRCode[];
  success?: boolean;
}

export interface HandlersRegisterRequest {
  company_name: string;
  email: string;
  /** @minLength 8 */
  password: string;
}

export interface HandlersResetPasswordRequest {
  /** @minLength 8 */
  new_password: string;
  token: string;
}

export interface ModelsAccount {
  company_name?: string;
  created_at?: string;
  email?: string;
  email_verified?: boolean;
  email_verified_at?: string;
  id?: string;
  is_active?: boolean;
  phone?: string;
  subscription_id?: string;
  /**
   * Subscription     *Subscription `json:"subscription,omitempty"` // TODO: Add when subscription domain is ready
   * Restaurants      []Restaurant  `json:"restaurants,omitempty"`  // TODO: Add when restaurant domain is ready
   */
  team_members?: ModelsTeamMember[];
  updated_at?: string;
}

export type ModelsFeedback = object;

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
  postal_code?: string;
  restaurant?: ModelsRestaurant;
  restaurant_id?: string;
  state?: string;
  updated_at?: string;
}

export interface ModelsPlanFeatures {
  advanced_analytics?: boolean;
  api_access?: boolean;
  custom_branding?: boolean;
  max_feedbacks_per_month?: number;
  max_locations_per_restaurant?: number;
  max_qr_codes_per_location?: number;
  max_restaurants?: number;
  max_team_members?: number;
  priority_support?: boolean;
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
  location?: ModelsLocation;
  location_id?: string;
  restaurant?: ModelsRestaurant;
  restaurant_id?: string;
  scans_count?: number;
  type?: ModelsQRCodeType;
  updated_at?: string;
}

export interface ModelsRestaurant {
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
  features?: ModelsPlanFeatures;
  id?: string;
  interval?: string;
  is_active?: boolean;
  name?: string;
  price?: number;
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
  role?: ModelsMemberRole;
  updated_at?: string;
  user?: ModelsUser;
  user_id?: string;
}

export interface ModelsUser {
  created_at?: string;
  email?: string;
  first_name?: string;
  id?: string;
  is_active?: boolean;
  last_name?: string;
  team_members?: ModelsTeamMember[];
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
 * @title LeCritique API
 * @version 1.0
 * @termsOfService http://swagger.io/terms/
 * @baseUrl //localhost:8080
 * @contact API Support <justdevelopitmx@proton.me>
 *
 * Restaurant feedback management system API
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  api = {
    /**
     * @description Check if the service is running
     *
     * @tags system
     * @name HealthList
     * @summary Health check
     * @request GET:/api/health
     */
    healthList: (params: RequestParams = {}) =>
      this.request<Record<string, any>, any>({
        path: `/api/health`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description Get detailed analytics data for a specific dish including ratings, feedback count, and recent feedback
     *
     * @tags analytics
     * @name V1AnalyticsDishesDetail
     * @summary Get dish analytics
     * @request GET:/api/v1/analytics/dishes/{dishId}
     * @secure
     */
    v1AnalyticsDishesDetail: (dishId: string, params: RequestParams = {}) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/analytics/dishes/${dishId}`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get comprehensive analytics data for a restaurant including ratings, feedback counts, and dish performance
     *
     * @tags analytics
     * @name V1AnalyticsRestaurantsDetail
     * @summary Get restaurant analytics
     * @request GET:/api/v1/analytics/restaurants/{restaurantId}
     * @secure
     */
    v1AnalyticsRestaurantsDetail: (
      restaurantId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/analytics/restaurants/${restaurantId}`,
        method: "GET",
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
     * @description Create a new restaurant owner account
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
     * @description Create a new dish for a restaurant
     *
     * @tags dishes
     * @name V1DishesCreate
     * @summary Create a new dish
     * @request POST:/api/v1/dishes
     * @secure
     */
    v1DishesCreate: (
      dish: HandlersCreateDishRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: GithubComLecritiqueApiInternalMenuModelsDish;
        },
        ResponseResponse
      >({
        path: `/api/v1/dishes`,
        method: "POST",
        body: dish,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get a specific dish by its ID
     *
     * @tags dishes
     * @name V1DishesDetail
     * @summary Get dish by ID
     * @request GET:/api/v1/dishes/{id}
     * @secure
     */
    v1DishesDetail: (id: string, params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: GithubComLecritiqueApiInternalMenuModelsDish;
        },
        ResponseResponse
      >({
        path: `/api/v1/dishes/${id}`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update a dish's information
     *
     * @tags dishes
     * @name V1DishesUpdate
     * @summary Update a dish
     * @request PUT:/api/v1/dishes/{id}
     * @secure
     */
    v1DishesUpdate: (
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
        path: `/api/v1/dishes/${id}`,
        method: "PUT",
        body: updates,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete a dish from the system
     *
     * @tags dishes
     * @name V1DishesDelete
     * @summary Delete a dish
     * @request DELETE:/api/v1/dishes/{id}
     * @secure
     */
    v1DishesDelete: (id: string, params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/dishes/${id}`,
        method: "DELETE",
        secure: true,
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
     * @description Submit customer feedback for a dish
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
     * @description Get questionnaire for a specific dish
     *
     * @tags public
     * @name V1PublicQuestionnaireDetail
     * @summary Get questionnaire
     * @request GET:/api/v1/public/questionnaire/{restaurantId}/{dishId}
     */
    v1PublicQuestionnaireDetail: (
      restaurantId: string,
      dishId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/public/questionnaire/${restaurantId}/${dishId}`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get public menu for a restaurant
     *
     * @tags public
     * @name V1PublicRestaurantMenuList
     * @summary Get restaurant menu
     * @request GET:/api/v1/public/restaurant/{id}/menu
     */
    v1PublicRestaurantMenuList: (id: string, params: RequestParams = {}) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/public/restaurant/${id}/menu`,
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
     * @description Get all restaurants for the authenticated account
     *
     * @tags restaurants
     * @name V1RestaurantsList
     * @summary Get all restaurants
     * @request GET:/api/v1/restaurants
     * @secure
     */
    v1RestaurantsList: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: ModelsRestaurant[];
        },
        ResponseResponse
      >({
        path: `/api/v1/restaurants`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Create a new restaurant for the authenticated account
     *
     * @tags restaurants
     * @name V1RestaurantsCreate
     * @summary Create a new restaurant
     * @request POST:/api/v1/restaurants
     * @secure
     */
    v1RestaurantsCreate: (
      request: HandlersCreateRestaurantRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: ModelsRestaurant;
        },
        ResponseResponse
      >({
        path: `/api/v1/restaurants`,
        method: "POST",
        body: request,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get a specific restaurant by its ID
     *
     * @tags restaurants
     * @name V1RestaurantsDetail
     * @summary Get restaurant by ID
     * @request GET:/api/v1/restaurants/{id}
     * @secure
     */
    v1RestaurantsDetail: (id: string, params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: ModelsRestaurant;
        },
        ResponseResponse
      >({
        path: `/api/v1/restaurants/${id}`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Update a restaurant's information
     *
     * @tags restaurants
     * @name V1RestaurantsUpdate
     * @summary Update restaurant
     * @request PUT:/api/v1/restaurants/{id}
     * @secure
     */
    v1RestaurantsUpdate: (
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
        path: `/api/v1/restaurants/${id}`,
        method: "PUT",
        body: updates,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Delete a restaurant from the system
     *
     * @tags restaurants
     * @name V1RestaurantsDelete
     * @summary Delete restaurant
     * @request DELETE:/api/v1/restaurants/{id}
     * @secure
     */
    v1RestaurantsDelete: (id: string, params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: Record<string, string>;
        },
        ResponseResponse
      >({
        path: `/api/v1/restaurants/${id}`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get feedback analytics and statistics for a restaurant
     *
     * @tags feedback
     * @name V1RestaurantsAnalyticsList
     * @summary Get feedback statistics
     * @request GET:/api/v1/restaurants/{restaurantId}/analytics
     * @secure
     */
    v1RestaurantsAnalyticsList: (
      restaurantId: string,
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/restaurants/${restaurantId}/analytics`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all dishes for a specific restaurant
     *
     * @tags dishes
     * @name V1RestaurantsDishesList
     * @summary Get dishes by restaurant
     * @request GET:/api/v1/restaurants/{restaurantId}/dishes
     * @secure
     */
    v1RestaurantsDishesList: (
      restaurantId: string,
      params: RequestParams = {},
    ) =>
      this.request<
        ResponseResponse & {
          data?: GithubComLecritiqueApiInternalMenuModelsDish[];
        },
        ResponseResponse
      >({
        path: `/api/v1/restaurants/${restaurantId}/dishes`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all feedback for a specific restaurant with pagination
     *
     * @tags feedback
     * @name V1RestaurantsFeedbackList
     * @summary Get restaurant feedback
     * @request GET:/api/v1/restaurants/{restaurantId}/feedback
     * @secure
     */
    v1RestaurantsFeedbackList: (
      restaurantId: string,
      query?: {
        /** Page number (default: 1) */
        page?: number;
        /** Items per page (default: 20, max: 100) */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<Record<string, any>, ResponseResponse>({
        path: `/api/v1/restaurants/${restaurantId}/feedback`,
        method: "GET",
        query: query,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all QR codes for a specific restaurant
     *
     * @tags qr-codes
     * @name V1RestaurantsQrCodesList
     * @summary Get QR codes by restaurant
     * @request GET:/api/v1/restaurants/{restaurantId}/qr-codes
     * @secure
     */
    v1RestaurantsQrCodesList: (
      restaurantId: string,
      params: RequestParams = {},
    ) =>
      this.request<HandlersQRCodeListResponse, ResponseResponse>({
        path: `/api/v1/restaurants/${restaurantId}/qr-codes`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Generate a new QR code for a restaurant
     *
     * @tags qr-codes
     * @name V1RestaurantsQrCodesCreate
     * @summary Generate QR code
     * @request POST:/api/v1/restaurants/{restaurantId}/qr-codes
     * @secure
     */
    v1RestaurantsQrCodesCreate: (
      restaurantId: string,
      qr_code: HandlersGenerateQRCodeRequest,
      params: RequestParams = {},
    ) =>
      this.request<HandlersGenerateQRCodeResponse, ResponseResponse>({
        path: `/api/v1/restaurants/${restaurantId}/qr-codes`,
        method: "POST",
        body: qr_code,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Check if the authenticated user can create more restaurants based on their subscription plan
     *
     * @tags subscription
     * @name V1UserCanCreateRestaurantList
     * @summary Check if user can create more restaurants
     * @request GET:/api/v1/user/can-create-restaurant
     * @secure
     */
    v1UserCanCreateRestaurantList: (params: RequestParams = {}) =>
      this.request<
        ResponseResponse & {
          data?: ServicesPermissionResponse;
        },
        ResponseResponse
      >({
        path: `/api/v1/user/can-create-restaurant`,
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
  };
}
