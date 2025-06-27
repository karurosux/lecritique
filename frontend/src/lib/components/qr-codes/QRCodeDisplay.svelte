<script lang="ts">
	import { onMount } from 'svelte';
	import { Modal, Button } from '$lib/components/ui';
	import { Download, Printer, Copy, Check } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import QRCode from 'qrcode';
	import { browser } from '$app/environment';

	let {
		qrCode,
		restaurantName,
		onclose
	}: {
		qrCode: any;
		restaurantName: string;
		onclose: () => void;
	} = $props();
	
	let canvas: HTMLCanvasElement;
	let qrDataUrl = $state('');
	let copied = $state(false);

	const qrUrl = browser ? `${window.location.origin}/qr/${qrCode.code}` : `/qr/${qrCode.code}`;

	onMount(async () => {
		try {
			// Generate QR code
			qrDataUrl = await QRCode.toDataURL(qrUrl, {
				width: 400,
				margin: 2,
				color: {
					dark: '#000000',
					light: '#FFFFFF'
				}
			});

			// Also render to canvas for better quality
			if (canvas) {
				await QRCode.toCanvas(canvas, qrUrl, {
					width: 400,
					margin: 2
				});
			}
		} catch (error) {
			console.error('Error generating QR code:', error);
			toast.error('Failed to generate QR code');
		}
	});

	function handleClose() {
		onclose();
	}

	async function handleCopyUrl() {
		try {
			await navigator.clipboard.writeText(qrUrl);
			copied = true;
			toast.success('URL copied to clipboard');
			setTimeout(() => {
				copied = false;
			}, 2000);
		} catch (error) {
			toast.error('Failed to copy URL');
		}
	}

	function handleDownload() {
		const link = document.createElement('a');
		link.download = `qr-code-${qrCode.label.toLowerCase().replace(/\s+/g, '-')}.png`;
		link.href = qrDataUrl;
		link.click();
	}

	function handlePrint() {
		if (!browser) return;
		
		const printWindow = window.open('', '_blank');
		if (!printWindow) {
			toast.error('Please allow popups to print');
			return;
		}

		const html = `
			<!DOCTYPE html>
			<html>
			<head>
				<title>QR Code - ${qrCode.label}</title>
				<style>
					body {
						font-family: system-ui, -apple-system, sans-serif;
						display: flex;
						flex-direction: column;
						align-items: center;
						justify-content: center;
						min-height: 100vh;
						margin: 0;
						padding: 20px;
					}
					.container {
						text-align: center;
						border: 2px solid #000;
						padding: 40px;
						border-radius: 8px;
					}
					h1 {
						margin: 0 0 10px 0;
						font-size: 24px;
					}
					h2 {
						margin: 0 0 20px 0;
						font-size: 20px;
						font-weight: normal;
					}
					.qr-code {
						margin: 20px 0;
					}
					.label {
						font-size: 18px;
						font-weight: bold;
						margin-top: 20px;
					}
					.instructions {
						margin-top: 20px;
						font-size: 14px;
						color: #666;
					}
					@media print {
						body {
							min-height: auto;
						}
					}
				</style>
			</head>
			<body>
				<div class="container">
					<h1>${restaurantName}</h1>
					<h2>Scan to leave feedback</h2>
					<img class="qr-code" src="${qrDataUrl}" alt="QR Code" width="300" height="300">
					<div class="label">${qrCode.label}</div>
					<div class="instructions">
						Scan this QR code with your phone camera<br>
						to share your dining experience
					</div>
				</div>
			</body>
			</html>
		`;

		printWindow.document.write(html);
		printWindow.document.close();
		printWindow.onload = () => {
			printWindow.print();
			printWindow.onafterprint = () => {
				printWindow.close();
			};
		};
	}
</script>

<Modal 
	isOpen={true} 
	title="QR Code: {qrCode.label}"
	size="xl"
	onclose={handleClose}
>
	<div class="space-y-4">
		<!-- QR Code Display -->
		<div class="flex justify-center p-8 bg-white rounded-lg border">
			<canvas bind:this={canvas} class="max-w-full"></canvas>
		</div>

		<!-- QR Code Info -->
		<div class="space-y-2 text-sm">
			<div class="flex justify-between">
				<span class="text-gray-500">Type:</span>
				<span class="capitalize">{qrCode.type}</span>
			</div>
			<div class="flex justify-between">
				<span class="text-gray-500">Code:</span>
				<span class="font-mono">{qrCode.code}</span>
			</div>
			<div class="flex justify-between">
				<span class="text-gray-500">Scans:</span>
				<span>{qrCode.scans_count || 0}</span>
			</div>
		</div>

		<!-- URL Copy -->
		<div class="flex items-center gap-2 p-3 bg-gray-50 rounded-lg">
			<code class="flex-1 text-sm truncate">{qrUrl}</code>
			<Button
				size="sm"
				variant="outline"
				onclick={handleCopyUrl}
			>
				{#if copied}
					<Check class="h-4 w-4" />
				{:else}
					<Copy class="h-4 w-4" />
				{/if}
			</Button>
		</div>

		<!-- Actions -->
		<div class="flex gap-2">
			<Button
				variant="outline"
				class="flex-1"
				onclick={handleDownload}
			>
				<Download class="mr-2 h-4 w-4" />
				Download
			</Button>
			<Button
				variant="outline"
				class="flex-1"
				onclick={handlePrint}
			>
				<Printer class="mr-2 h-4 w-4" />
				Print
			</Button>
		</div>

		<!-- Instructions -->
		<div class="rounded-lg bg-gray-50 p-4 text-sm">
			<h4 class="font-medium mb-2">Usage Instructions</h4>
			<ol class="space-y-1 text-gray-600">
				<li>1. Download or print this QR code</li>
				<li>2. Place it at {qrCode.label}</li>
				<li>3. Customers scan with their phone camera</li>
				<li>4. They'll be directed to leave feedback</li>
			</ol>
		</div>
	</div>
</Modal>