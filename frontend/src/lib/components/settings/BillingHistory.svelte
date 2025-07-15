<script lang="ts">
  import { APP_CONFIG } from '$lib/constants/config';
  
  interface Invoice {
    id: string;
    plan: string;
    amount: number;
    date: string;
    status: 'paid' | 'pending' | 'failed';
  }

  interface Props {
    invoices?: Invoice[];
    onDownload?: (invoiceId: string) => void;
  }

  let { 
    invoices = [
      { id: '1', plan: 'Professional Plan', amount: 79, date: 'December 15, 2024', status: 'paid' },
      { id: '2', plan: 'Professional Plan', amount: 79, date: 'November 15, 2024', status: 'paid' },
      { id: '3', plan: 'Professional Plan', amount: 79, date: 'October 15, 2024', status: 'paid' },
      { id: '4', plan: 'Professional Plan', amount: 79, date: 'September 15, 2024', status: 'paid' },
      { id: '5', plan: 'Starter Plan', amount: 29, date: 'August 15, 2024', status: 'paid' }
    ] as Invoice[],
    onDownload
  }: Props = $props();

  function handleDownload(invoiceId: string) {
    onDownload?.(invoiceId);
  }
</script>

<div>
  <div class="mb-8">
    <h2 class="text-2xl font-bold text-gray-900">Billing History</h2>
    <p class="mt-1 text-sm text-gray-600">View and download your past invoices</p>
  </div>
  
  <div class="space-y-3">
    {#each invoices as invoice}
      <div class="flex items-center justify-between p-4 rounded-lg border border-gray-200 hover:border-gray-300 transition-colors">
        <div>
          <p class="font-medium text-gray-900">{invoice.plan}</p>
          <p class="text-sm text-gray-600">{invoice.date}</p>
        </div>
        <div class="text-right">
          <p class="font-semibold text-gray-900">${invoice.amount.toFixed(2)}</p>
          <button 
            class="text-sm text-blue-600 hover:text-blue-700 font-medium"
            onclick={() => handleDownload(invoice.id)}
          >
            Download
          </button>
        </div>
      </div>
    {/each}
  </div>

  <div class="mt-6 pt-6 border-t">
    <p class="text-sm text-gray-600">
      All invoices include tax information and can be used for accounting purposes. 
      For billing inquiries, please contact {APP_CONFIG.emails.billing}
    </p>
  </div>
</div>