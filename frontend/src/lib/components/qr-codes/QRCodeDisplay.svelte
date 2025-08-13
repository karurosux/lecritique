<script lang="ts">
	import { onMount } from 'svelte';
	import { Modal, Button } from '$lib/components/ui';
	import { Download, Printer, Copy, Check } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import QRCode from 'qrcode';
	import { browser } from '$app/environment';

	let {
		qrCode,
		organizationName,
		clickOrigin = null,
		onclose
	}: {
		qrCode: any;
		organizationName: string;
		clickOrigin?: { x: number; y: number } | null;
		onclose: () => void;
	} = $props();

	let canvas: HTMLCanvasElement | undefined;
	let qrDataUrl = $state('');
	let copied = $state(false);
	let showPrintOptions = $state(false);

	const qrUrl = browser ? `${window.location.origin}/qr/${qrCode.code}` : `/qr/${qrCode.code}`;

	onMount(async () => {
		try {
			qrDataUrl = await QRCode.toDataURL(qrUrl, {
				width: 400,
				margin: 2,
				color: {
					dark: '#000000',
					light: '#FFFFFF'
				}
			});

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

		printWindow.document.write('<!DOCTYPE html>');
		printWindow.document.write('<html><head>');
		printWindow.document.write('<title>QR Codes - ' + qrCode.label + '</title>');

		printWindow.document.write('<sty' + 'le>');
		printWindow.document.write('* { box-sizing: border-box; }');
		printWindow.document.write(
			'body { font-family: system-ui, -apple-system, sans-serif; margin: 0; padding: 20px; }'
		);
		printWindow.document.write('.page { max-width: 8.5in; margin: 0 auto; }');
		printWindow.document.write(
			'.instructions-header { text-align: center; margin-bottom: 30px; padding: 20px; background: #f5f5f5; border-radius: 8px; }'
		);
		printWindow.document.write('.instructions-header h1 { margin: 0 0 10px 0; font-size: 24px; }');
		printWindow.document.write(
			'.instructions-header p { margin: 5px 0; color: #666; font-size: 14px; }'
		);
		printWindow.document.write('.section { margin: 20px 0; }');
		printWindow.document.write(
			'.qr-group { display: flex; justify-content: center; align-items: flex-start; flex-wrap: wrap; gap: 0; }'
		);
		printWindow.document.write(
			'.qr-item { text-align: center; border: 1px dashed rgba(100, 100, 100, 0.5); padding: 15px; position: relative; }'
		);
		printWindow.document.write(
			'.qr-item + .qr-item { border-left: 1px dashed rgba(100, 100, 100, 0.5); }'
		);
		printWindow.document.write(
			'.qr-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0; justify-content: center; }'
		);
		printWindow.document.write(
			'.qr-grid .qr-item:nth-child(1), .qr-grid .qr-item:nth-child(2) { border-bottom: 1px dashed rgba(100, 100, 100, 0.5); }'
		);
		printWindow.document.write(
			'.qr-grid .qr-item:nth-child(2), .qr-grid .qr-item:nth-child(4) { border-left: 1px dashed rgba(100, 100, 100, 0.5); }'
		);
		printWindow.document.write('.qr-item img { display: block; margin: 0 auto; }');
		printWindow.document.write(
			'.qr-item .logo-text { font-size: 20px; font-weight: bold; color: #4F46E5; margin-bottom: 10px; }'
		);
		printWindow.document.write('.qr-item h3 { margin: 10px 0 5px 0; font-size: 14px; }');
		printWindow.document.write(
			'.qr-item .label { font-size: 12px; font-weight: bold; margin: 5px 0; }'
		);
		printWindow.document.write('.qr-item .use-case { font-size: 10px; color: #666; }');
		printWindow.document.write(
			'.cut-line { border: 0; height: 0; border-top: 1px dashed #999; margin: 0; }'
		);
		printWindow.document.write('@media print { body { padding: 0; } @page { margin: 0.5in; } }');
		printWindow.document.write('</sty' + 'le>');

		printWindow.document.write('</head><body>');
		printWindow.document.write('<div class="page">');

		printWindow.document.write('<div class="instructions-header">');
		printWindow.document.write('<h1>' + organizationName + ' - QR Code Sheet</h1>');
		printWindow.document.write(
			'<p>Cut along the dashed lines to separate individual QR code cards</p>'
		);
		printWindow.document.write(
			'<p>All QR codes lead to the same feedback form for: <strong>' +
				qrCode.label +
				'</strong></p>'
		);
		printWindow.document.write('</div>');

		printWindow.document.write('<div class="section">');
		printWindow.document.write('<div class="qr-group">');
		printWindow.document.write('<div class="qr-item">');
		printWindow.document.write('<div class="logo-text">Kyooar</div>');
		printWindow.document.write(
			'<img src="' + qrDataUrl + '" alt="QR Code" width="200" height="200">'
		);
		printWindow.document.write('<h3>' + organizationName + '</h3>');
		printWindow.document.write('<div class="label">' + qrCode.label + '</div>');
		printWindow.document.write('<div class="use-case">Scan to leave feedback</div>');
		printWindow.document.write('</div>');
		printWindow.document.write('</div>');
		printWindow.document.write('</div>');

		printWindow.document.write('<div class="section">');
		printWindow.document.write('<div class="qr-group">');
		for (let i = 0; i < 2; i++) {
			printWindow.document.write('<div class="qr-item">');
			printWindow.document.write('<div class="logo-text">Kyooar</div>');
			printWindow.document.write(
				'<img src="' + qrDataUrl + '" alt="QR Code" width="120" height="120">'
			);
			printWindow.document.write('<div class="label">' + qrCode.label + '</div>');
			printWindow.document.write('<div class="use-case">Share your experience</div>');
			printWindow.document.write('</div>');
		}
		printWindow.document.write('</div>');
		printWindow.document.write('</div>');

		printWindow.document.write('<div class="section">');
		printWindow.document.write('<div class="qr-grid">');
		for (let i = 0; i < 4; i++) {
			printWindow.document.write('<div class="qr-item">');
			printWindow.document.write(
				'<img src="' + qrDataUrl + '" alt="QR Code" width="80" height="80">'
			);
			printWindow.document.write(
				'<div class="use-case" style="font-size: 9px;">Feedback: ' + qrCode.label + '</div>'
			);
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

	function handlePrintTentCard() {
		if (!browser) return;

		const printWindow = window.open('', '_blank');
		if (!printWindow) {
			toast.error('Please allow popups to print');
			return;
		}

		printWindow.document.write('<!DOCTYPE html>');
		printWindow.document.write('<html><head>');
		printWindow.document.write('<title>QR Display Stand - ' + qrCode.label + '</title>');

		printWindow.document.write('<sty' + 'le>');
		printWindow.document.write('* { margin: 0; padding: 0; box-sizing: border-box; }');
		printWindow.document.write(
			'body { font-family: "Inter", sans-serif; background: white; padding: 20mm; }'
		);
		printWindow.document.write(
			'.page-title { text-align: center; font-size: 24px; font-weight: bold; margin-bottom: 20px; color: #333; }'
		);
		printWindow.document.write(
			'.legend { background: #f9f9f9; border: 1px solid #ddd; padding: 15px; margin-bottom: 20px; border-radius: 8px; }'
		);
		printWindow.document.write('.legend h3 { margin-bottom: 10px; color: #333; }');
		printWindow.document.write(
			'.legend-items { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }'
		);
		printWindow.document.write('.legend-item { display: flex; align-items: center; gap: 10px; }');
		printWindow.document.write('.legend-line { width: 30px; height: 3px; }');
		printWindow.document.write('.cut-line-legend { background: #000; }');
		printWindow.document.write(
			'.valley-fold-legend { background: repeating-linear-gradient(to right, #0066cc 0, #0066cc 5px, transparent 5px, transparent 10px); height: 2px; }'
		);
		printWindow.document.write('.svg-container { text-align: center; margin: 20px 0; }');
		printWindow.document.write(
			'.template-svg { max-width: 100%; height: auto; border: 1px solid #ccc; }'
		);
		printWindow.document.write(
			'.instructions { background: #f5f5f5; padding: 20px; border-radius: 8px; margin-top: 30px; }'
		);
		printWindow.document.write('.instructions h3 { color: #555; margin: 15px 0 8px 0; }');
		printWindow.document.write('.instructions ol { margin-left: 20px; margin-bottom: 15px; }');
		printWindow.document.write('.instructions li { margin-bottom: 5px; line-height: 1.4; }');
		printWindow.document.write('@page { size: A4; margin: 20mm; }');
		printWindow.document.write('@media print { body { background: white; padding: 0; } }');
		printWindow.document.write('</sty' + 'le>');

		printWindow.document.write('</head><body>');

		printWindow.document.write('<h1 class="page-title">QR Display Stand Template</h1>');

		printWindow.document.write('<div class="legend">');
		printWindow.document.write('<h3>Legend</h3>');
		printWindow.document.write('<div class="legend-items">');
		printWindow.document.write(
			'<div class="legend-item"><div class="legend-line cut-line-legend"></div><span>Cut line - Cut with scissors</span></div>'
		);
		printWindow.document.write(
			'<div class="legend-item"><div class="legend-line valley-fold-legend"></div><span>Valley fold - Fold toward you</span></div>'
		);
		printWindow.document.write('</div></div>');

		printWindow.document.write('<div class="svg-container">');
		// printWindow.document.write(
		//   '<svg class="template-svg" viewBox="0 0 600 400" xmlns="http:
		// );

		printWindow.document.write(
			'<rect x="100" y="50" width="180" height="220" fill="white" stroke="black" stroke-width="3"/>'
		);

		printWindow.document.write(
			'<text x="190" y="90" text-anchor="middle" font-size="24" font-weight="bold" fill="#4F46E5">Kyooar</text>'
		);

		printWindow.document.write(
			'<text x="190" y="120" text-anchor="middle" font-size="16" font-weight="bold" fill="#333">' +
				organizationName +
				'</text>'
		);

		printWindow.document.write(
			'<image x="125" y="130" width="130" height="130" href="' + qrDataUrl + '"/>'
		);

		printWindow.document.write(
			'<text x="190" y="280" text-anchor="middle" font-size="14" font-weight="bold" fill="#000">' +
				qrCode.label +
				'</text>'
		);

		if (qrCode.location) {
			printWindow.document.write(
				'<text x="190" y="295" text-anchor="middle" font-size="11" fill="#666">' +
					qrCode.location +
					'</text>'
			);
		}

		printWindow.document.write(
			'<rect x="280" y="100" width="100" height="120" fill="#f5f5f5" stroke="black" stroke-width="2"/>'
		);
		printWindow.document.write(
			'<text x="330" y="165" text-anchor="middle" font-size="10" fill="#999">BACK SUPPORT</text>'
		);

		printWindow.document.write(
			'<rect x="100" y="270" width="180" height="60" fill="#e8e8e8" stroke="black" stroke-width="2"/>'
		);
		printWindow.document.write(
			'<text x="190" y="305" text-anchor="middle" font-size="10" fill="#666">BASE</text>'
		);

		printWindow.document.write(
			'<line x1="280" y1="100" x2="280" y2="220" stroke="#0066cc" stroke-width="2" stroke-dasharray="6,6"/>'
		);
		printWindow.document.write(
			'<line x1="100" y1="270" x2="280" y2="270" stroke="#0066cc" stroke-width="2" stroke-dasharray="6,6"/>'
		);

		printWindow.document.write(
			'<text x="420" y="100" font-size="12" font-weight="bold">Simple Assembly:</text>'
		);
		printWindow.document.write(
			'<text x="420" y="120" font-size="10">1. Cut along black lines</text>'
		);
		printWindow.document.write(
			'<text x="420" y="135" font-size="10">2. Fold back support 90°</text>'
		);
		printWindow.document.write(
			'<text x="420" y="150" font-size="10">3. Fold base under 90°</text>'
		);
		printWindow.document.write(
			'<text x="420" y="165" font-size="10">4. Creates inclined stand</text>'
		);

		printWindow.document.write(
			'<rect x="440" y="180" width="60" height="50" fill="white" stroke="#333" stroke-width="2"/>'
		);
		printWindow.document.write(
			'<polygon points="440,180 460,170 460,220 440,230" fill="#f0f0f0" stroke="#333" stroke-width="1"/>'
		);
		printWindow.document.write(
			'<rect x="440" y="230" width="60" height="10" fill="#e0e0e0" stroke="#333" stroke-width="1"/>'
		);
		printWindow.document.write(
			'<text x="420" y="260" font-size="9" fill="#666">L-shaped stand</text>'
		);

		printWindow.document.write('</svg></div>');

		printWindow.document.write('<div class="instructions">');
		printWindow.document.write('<h3>L-Shaped Stand Assembly:</h3>');
		printWindow.document.write('<ol>');
		printWindow.document.write('<li><strong>Cut:</strong> Cut along all solid black lines</li>');
		printWindow.document.write(
			'<li><strong>Fold back support:</strong> Fold the side panel 90° backward to create the support</li>'
		);
		printWindow.document.write(
			'<li><strong>Fold base:</strong> Fold the bottom panel 90° under to create the base</li>'
		);
		printWindow.document.write(
			'<li><strong>Done:</strong> The L-shape creates a stable inclined stand</li>'
		);
		printWindow.document.write('</ol>');

		printWindow.document.write('<p style="margin-top: 15px; font-style: italic; color: #666;">');
		printWindow.document.write(
			'Simple L-shaped design: front displays QR code, back provides support, base provides stability and incline.'
		);
		printWindow.document.write('</p>');
		printWindow.document.write('</div>');

		printWindow.document.write('</body></html>');

		printWindow.document.close();
		printWindow.onload = () => {
			printWindow.print();
			printWindow.onafterprint = () => {
				printWindow.close();
			};
		};
		showPrintOptions = false;
	}
</script>

<Modal isOpen={true} title="QR Code: {qrCode.label}" {clickOrigin} size="xl" onclose={handleClose}>
	<div class="space-y-4">
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

		<div class="flex items-center gap-2 rounded-lg bg-gray-50 p-3">
			<code class="flex-1 truncate text-sm">{qrUrl}</code>
			<Button size="sm" variant="outline" onclick={handleCopyUrl}>
				{#if copied}
					<Check class="h-4 w-4" />
				{:else}
					<Copy class="h-4 w-4" />
				{/if}
			</Button>
		</div>

		<div class="flex gap-2">
			<Button variant="outline" class="flex-1" onclick={handleDownload}>
				<Download class="mr-2 h-4 w-4" />
				Download
			</Button>
			<div class="relative flex-1">
				<Button
					variant="outline"
					class="w-full"
					onclick={(e) => {
						e.stopPropagation();
						showPrintOptions = !showPrintOptions;
					}}
				>
					<Printer class="mr-2 h-4 w-4" />
					Print
				</Button>

				{#if showPrintOptions}
					<div class="fixed inset-0 z-10" onclick={() => (showPrintOptions = false)}></div>
					<div class="absolute top-full z-20 mt-2 w-full rounded-lg border bg-white shadow-lg">
						<button
							class="w-full rounded-lg px-4 py-2 text-left text-sm hover:bg-gray-50"
							onclick={(e) => {
								e.stopPropagation();
								handlePrint();
							}}
						>
							<div class="font-medium">Print Cards</div>
							<div class="text-xs text-gray-500">Multiple sizes with cut guides</div>
						</button>
					</div>
				{/if}
			</div>
		</div>

		<div class="rounded-lg bg-gray-50 p-4 text-sm">
			<h4 class="mb-2 font-medium">Usage Instructions</h4>
			<ol class="space-y-1 text-gray-600">
				<li>1. Download or print this QR code</li>
				<li>2. Place it at {qrCode.label}</li>
				<li>3. Customers scan with their phone camera</li>
				<li>4. They'll be directed to leave feedback</li>
			</ol>
		</div>
	</div>
</Modal>
