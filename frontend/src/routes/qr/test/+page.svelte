<script lang="ts">
  import { Card, Button } from '$lib/components/ui';
  import { goto } from '$app/navigation';

  function testQRCode() {
    goto('/qr/DEMO001');
  }

  function testRestaurantMenu() {
    // Using a sample restaurant ID - in real usage this would come from QR validation
    goto('/restaurant/27178b57-2dd4-4311-aea5-f239668e79ae/menu?qr=DEMO001');
  }
</script>

<svelte:head>
  <title>QR Code Test - LeCritique</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 py-8 px-4">
  <div class="max-w-2xl mx-auto">
    <Card>
      <div class="text-center py-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-4">QR Code System Test</h1>
        <p class="text-gray-600 mb-6">Test the QR code validation and restaurant menu pages</p>
        
        <div class="space-y-4">
          <Button variant="primary" on:click={testQRCode} class="w-full">
            Test QR Code: DEMO001
          </Button>
          
          <Button variant="secondary" on:click={testRestaurantMenu} class="w-full">
            Test Restaurant Menu (Direct)
          </Button>
        </div>
        
        <div class="mt-8 text-sm text-gray-500">
          <p><strong>Note:</strong> Make sure the backend is running and seeded with:</p>
          <code class="bg-gray-100 px-2 py-1 rounded">make seed</code>
        </div>
      </div>
    </Card>
  </div>
</div>