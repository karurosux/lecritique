<script lang="ts">
  import { onMount } from 'svelte';
  import QRCodeLib from 'qrcode';

  let {
    data = '',
    size = 128,
    margin = 0,
    errorCorrectionLevel = 'M',
    color = {
      dark: '#000000',
      light: '#FFFFFF',
    },
  }: {
    data?: string;
    size?: number;
    margin?: number;
    errorCorrectionLevel?: 'L' | 'M' | 'Q' | 'H';
    color?: {
      dark: string;
      light: string;
    };
  } = $props();

  let canvas: HTMLCanvasElement;

  onMount(() => {
    if (canvas && data) {
      QRCodeLib.toCanvas(canvas, data, {
        width: size,
        margin: margin,
        errorCorrectionLevel: errorCorrectionLevel,
        color: color,
      });
    }
  });

  $effect(() => {
    if (canvas && data) {
      QRCodeLib.toCanvas(canvas, data, {
        width: size,
        margin: margin,
        errorCorrectionLevel: errorCorrectionLevel,
        color: color,
      });
    }
  });
</script>

<canvas bind:this={canvas} width={size} height={size}></canvas>
