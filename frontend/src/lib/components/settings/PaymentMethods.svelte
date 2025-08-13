<script lang="ts">
  import { Button } from '$lib/components/ui';
  import { CreditCard, Building, AlertCircle } from 'lucide-svelte';

  interface PaymentMethod {
    id: string;
    type: 'card' | 'bank';
    brand?: string;
    last4: string;
    expiryMonth?: number;
    expiryYear?: number;
    isDefault: boolean;
  }

  interface BillingAddress {
    company: string;
    street1: string;
    street2?: string;
    city: string;
    state: string;
    postalCode: string;
    country: string;
  }

  interface Props {
    paymentMethods?: PaymentMethod[];
    billingAddress?: BillingAddress;
    onEditPayment?: (id: string) => void;
    onAddPayment?: (type: 'card' | 'bank') => void;
    onUpdateAddress?: () => void;
  }

  let {
    paymentMethods = [
      {
        id: '1',
        type: 'card',
        brand: 'VISA',
        last4: '4242',
        expiryMonth: 12,
        expiryYear: 2025,
        isDefault: true,
      },
    ] as PaymentMethod[],
    billingAddress = {
      company: 'Acme Organization Group',
      street1: '123 Main Street',
      street2: 'Suite 100',
      city: 'San Francisco',
      state: 'CA',
      postalCode: '94105',
      country: 'United States',
    } as BillingAddress,
    onEditPayment,
    onAddPayment,
    onUpdateAddress,
  }: Props = $props();
</script>

<div>
  <div class="mb-8">
    <h2 class="text-2xl font-bold text-gray-900">Payment Methods</h2>
    <p class="mt-1 text-sm text-gray-600">
      Manage your payment methods and billing address
    </p>
  </div>

  
  <div class="space-y-4">
    <h3 class="text-lg font-medium text-gray-900">Current Payment Method</h3>

    {#each paymentMethods as method}
      <div
        class="border-2 rounded-xl p-6 {method.isDefault
          ? 'border-blue-200 bg-blue-50'
          : 'border-gray-200'}">
        <div class="flex items-start justify-between">
          <div class="flex items-center space-x-4">
            {#if method.type === 'card'}
              <div
                class="h-12 w-20 bg-gradient-to-r from-blue-600 to-blue-400 rounded-lg flex items-center justify-center text-white font-bold">
                {method.brand}
              </div>
              <div>
                <p class="font-semibold text-gray-900">
                  •••• •••• •••• {method.last4}
                </p>
                <p class="text-sm text-gray-600">
                  Expires {method.expiryMonth}/{method.expiryYear}
                </p>
                {#if method.isDefault}
                  <p class="text-xs text-blue-600 mt-1">
                    Default payment method
                  </p>
                {/if}
              </div>
            {:else}
              <div
                class="h-12 w-20 bg-gradient-to-r from-gray-600 to-gray-400 rounded-lg flex items-center justify-center text-white">
                <Building class="h-6 w-6" />
              </div>
              <div>
                <p class="font-semibold text-gray-900">
                  Bank Account •••• {method.last4}
                </p>
                {#if method.isDefault}
                  <p class="text-xs text-blue-600 mt-1">
                    Default payment method
                  </p>
                {/if}
              </div>
            {/if}
          </div>
          <Button
            variant="outline"
            size="sm"
            onclick={() => onEditPayment?.(method.id)}>
            Edit
          </Button>
        </div>
      </div>
    {/each}

    
    <div class="mt-8">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Add Payment Method</h3>

      <div class="grid gap-4 md:grid-cols-2">
        <button
          class="p-6 border-2 border-dashed border-gray-300 rounded-xl hover:border-gray-400 transition-colors text-center group"
          onclick={() => onAddPayment?.('card')}>
          <CreditCard
            class="h-8 w-8 text-gray-400 mx-auto mb-2 group-hover:text-gray-500" />
          <p class="text-sm font-medium text-gray-900">
            Add Credit or Debit Card
          </p>
        </button>

        <button
          class="p-6 border-2 border-dashed border-gray-300 rounded-xl hover:border-gray-400 transition-colors text-center group"
          onclick={() => onAddPayment?.('bank')}>
          <Building
            class="h-8 w-8 text-gray-400 mx-auto mb-2 group-hover:text-gray-500" />
          <p class="text-sm font-medium text-gray-900">Link Bank Account</p>
        </button>
      </div>
    </div>

    
    <div class="mt-8 border-t pt-8">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Billing Address</h3>

      <div class="bg-gray-50 rounded-lg p-4">
        <p class="font-medium text-gray-900">{billingAddress.company}</p>
        <p class="text-sm text-gray-600 mt-1">
          {billingAddress.street1}<br />
          {#if billingAddress.street2}
            {billingAddress.street2}<br />
          {/if}
          {billingAddress.city}, {billingAddress.state}
          {billingAddress.postalCode}<br />
          {billingAddress.country}
        </p>
        <Button
          variant="outline"
          size="sm"
          class="mt-4"
          onclick={onUpdateAddress}>
          Update Address
        </Button>
      </div>
    </div>

    
    <div class="mt-8 border-t pt-8">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Security</h3>

      <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
        <div class="flex">
          <AlertCircle class="h-5 w-5 text-yellow-400 flex-shrink-0" />
          <div class="ml-3">
            <p class="text-sm text-yellow-800">
              Your payment information is encrypted and stored securely. We
              never store your full card details.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
