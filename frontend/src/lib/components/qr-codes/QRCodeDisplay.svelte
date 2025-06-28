<script lang="ts">
	import { onMount } from 'svelte';
	import { Modal, Button } from '$lib/components/ui';
	import { Download, Printer, Copy, Check } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import QRCode from 'qrcode';
	import { browser } from '$app/environment';
	import logoImage from '../../../assets/logo.png';

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

		// Build the print document
		printWindow.document.write('<!DOCTYPE html>');
		printWindow.document.write('<html><head>');
		printWindow.document.write('<title>QR Codes - ' + qrCode.label + '</title>');
		
		// Write styles directly to avoid preprocessor issues
		printWindow.document.write('<sty' + 'le>');
		printWindow.document.write('* { box-sizing: border-box; }');
		printWindow.document.write('body { font-family: system-ui, -apple-system, sans-serif; margin: 0; padding: 20px; }');
		printWindow.document.write('.page { max-width: 8.5in; margin: 0 auto; }');
		printWindow.document.write('.instructions-header { text-align: center; margin-bottom: 30px; padding: 20px; background: #f5f5f5; border-radius: 8px; }');
		printWindow.document.write('.instructions-header h1 { margin: 0 0 10px 0; font-size: 24px; }');
		printWindow.document.write('.instructions-header p { margin: 5px 0; color: #666; font-size: 14px; }');
		printWindow.document.write('.section { margin: 20px 0; }');
		printWindow.document.write('.qr-group { display: flex; justify-content: center; align-items: flex-start; flex-wrap: wrap; gap: 0; }');
		printWindow.document.write('.qr-item { text-align: center; border: 1px dashed rgba(100, 100, 100, 0.5); padding: 15px; position: relative; }');
		printWindow.document.write('.qr-item + .qr-item { border-left: 1px dashed rgba(100, 100, 100, 0.5); }');
		printWindow.document.write('.qr-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0; justify-content: center; }');
		printWindow.document.write('.qr-grid .qr-item:nth-child(1), .qr-grid .qr-item:nth-child(2) { border-bottom: 1px dashed rgba(100, 100, 100, 0.5); }');
		printWindow.document.write('.qr-grid .qr-item:nth-child(2), .qr-grid .qr-item:nth-child(4) { border-left: 1px dashed rgba(100, 100, 100, 0.5); }');
		printWindow.document.write('.qr-item img { display: block; margin: 0 auto; }');
		printWindow.document.write('.qr-item .logo { height: 30px; margin-bottom: 10px; }');
		printWindow.document.write('.qr-item h3 { margin: 10px 0 5px 0; font-size: 14px; }');
		printWindow.document.write('.qr-item .label { font-size: 12px; font-weight: bold; margin: 5px 0; }');
		printWindow.document.write('.qr-item .use-case { font-size: 10px; color: #666; }');
		printWindow.document.write('.cut-line { border: 0; height: 0; border-top: 1px dashed #999; margin: 0; }');
		printWindow.document.write('@media print { body { padding: 0; } @page { margin: 0.5in; } }');
		printWindow.document.write('</sty' + 'le>');
		
		printWindow.document.write('</head><body>');
		printWindow.document.write('<div class="page">');
		
		// Instructions header
		printWindow.document.write('<div class="instructions-header">');
		printWindow.document.write('<h1>' + restaurantName + ' - QR Code Sheet</h1>');
		printWindow.document.write('<p>Cut along the dashed lines to separate individual QR code cards</p>');
		printWindow.document.write('<p>All QR codes lead to the same feedback form for: <strong>' + qrCode.label + '</strong></p>');
		printWindow.document.write('</div>');
		
		// Large format section
		printWindow.document.write('<div class="section">');
		printWindow.document.write('<div class="qr-group">');
		printWindow.document.write('<div class="qr-item">');
		printWindow.document.write('<img class="logo" src="' + logoImage + '" alt="LeCritique Logo">');
		printWindow.document.write('<img src="' + qrDataUrl + '" alt="QR Code" width="200" height="200">');
		printWindow.document.write('<h3>' + restaurantName + '</h3>');
		printWindow.document.write('<div class="label">' + qrCode.label + '</div>');
		printWindow.document.write('<div class="use-case">Scan to leave feedback</div>');
		printWindow.document.write('</div>');
		printWindow.document.write('</div>');
		printWindow.document.write('</div>');
		
		// Medium format section
		printWindow.document.write('<div class="section">');
		printWindow.document.write('<div class="qr-group">');
		for (let i = 0; i < 2; i++) {
			printWindow.document.write('<div class="qr-item">');
			printWindow.document.write('<img class="logo" src="' + logoImage + '" alt="LeCritique Logo">');
			printWindow.document.write('<img src="' + qrDataUrl + '" alt="QR Code" width="120" height="120">');
			printWindow.document.write('<div class="label">' + qrCode.label + '</div>');
			printWindow.document.write('<div class="use-case">Share your experience</div>');
			printWindow.document.write('</div>');
		}
		printWindow.document.write('</div>');
		printWindow.document.write('</div>');
		
		// Small format section (2x2 grid)
		printWindow.document.write('<div class="section">');
		printWindow.document.write('<div class="qr-grid">');
		for (let i = 0; i < 4; i++) {
			printWindow.document.write('<div class="qr-item">');
			printWindow.document.write('<img src="' + qrDataUrl + '" alt="QR Code" width="80" height="80">');
			printWindow.document.write('<div class="use-case" style="font-size: 9px;">Feedback: ' + qrCode.label + '</div>');
			printWindow.document.write('</div>');
		}
		printWindow.document.write('</div>');
		printWindow.document.write('</div>');
		
		printWindow.document.write('</div>');
		printWindow.document.write('</body></html>');
		
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