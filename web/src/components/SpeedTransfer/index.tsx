export type SpeedTransferProps = {
  bytes: number
  decimals?: number
}

export const SpeedTransfer: (props: SpeedTransferProps) => (string | string) = (props: SpeedTransferProps) => {
  if (props.decimals === undefined || props.decimals === 0) {
    props.decimals = 2
  }
  if (props.bytes === 0) return '0 B';

  const k = 1024;
  const dm = props.decimals < 0 ? 0 : props.decimals;
  const sizes = ["B", 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(props.bytes) / Math.log(k));

  return parseFloat((props.bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}
