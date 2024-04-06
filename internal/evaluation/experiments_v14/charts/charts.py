import os
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

def validate_file(filepath, kind):
  """
  Checks if the file exists and the first line matches the expected header.

  Args:
      filepath (str): Path to the file.
      kind: File kind (monitor or results).

  Returns:
      bool: True if the file is valid, False otherwise.
  """
  try:
    with open(filepath, 'r') as f:
      first_line = f.readline().strip()
      if kind == "results":
        return first_line == "dateTime;info;sequential;response_time"
      else:
        return first_line == "dateTime;container_name;container_status;used_memory(MB);available_memory(MB);memory_usage(%);cpu_delta;system_cpu_delta;number_cpus;cpu_usage(%);total_cpu_usage;pre_total_cpu_usage"
  except FileNotFoundError:
    print(f"File not found: {filepath}")
    return False

def get_last_log_file(directory, app, kind):
  """
  Identifica e retorna o último arquivo de log em uma pasta.

  Args:
    directory: Caminho para a pasta que contém os arquivos de log.
    kind: Tipo de arquivo de log (monitor ou results).

  Returns:
    Caminho para o último arquivo de log.
  """
  files = os.listdir(directory)
  kindFilter = f".{kind}." + "csv" if kind == "monitor" else "txt"
  logFiles = [f for f in files if f.endswith(kindFilter) and app in f]
  if not logFiles:
    return None
  # return os.path.join(directory, max(logFiles))
  # if validate_file(os.path.join(directory, max(log_files)), kind):

  logFiles.sort(reverse=True)

  for logFile in logFiles:
    if validate_file(os.path.join(directory, logFile), kind):
      return os.path.join(directory, max(logFiles))
  return None


def calculate_duration(df):
  """
  Calcula a duração do experimento em segundos.

  Args:
    df: DataFrame do Pandas com os dados do experimento.

  Returns:
    Duração do experimento em segundos.
  """
  start_time = pd.to_datetime(df["dateTime"].iloc[0])
  end_time = pd.to_datetime(df["dateTime"].iloc[-1])
  return (end_time - start_time).total_seconds()

def read_monitor_data(file_path):
  """
  Lê os dados do arquivo de log e os armazena em um DataFrame do Pandas.

  Args:
    file_path: Caminho para o arquivo de log.

  Returns:
    DataFrame do Pandas com os dados do experimento.
  """
  df = pd.read_csv(file_path, delimiter=";")
  df["dateTime"] = pd.to_datetime(df["dateTime"], format="mixed")
  df['duration'] = (df['dateTime'] - df['dateTime'].iloc[0]).dt.total_seconds()
  return df

def read_results_data(file_path):
  """
  Lê os dados do arquivo de log e os armazena em um DataFrame do Pandas.

  Args:
    file_path: Caminho para o arquivo de log.

  Returns:
    DataFrame do Pandas com os dados do experimento.
  """
  header = ["dateTime", "info", "sequential", "response_time"]
  # print(file_path)

  df = pd.read_csv(file_path, header=None, delimiter=";", skiprows=101, names=header, na_values=["-"]) # 101 = 100 Warm-up requests + header


  # # Definindo a largura da coluna fixa
  # width_fixa = 19
  # # Lendo a coluna fixa com `pd.read_fwf`
  # df_fixa = pd.read_fwf(file_path, widths=[width_fixa], skiprows=1, header=None)
  # # Lendo as demais colunas com `pd.read_csv`
  # df_variavel = pd.read_csv(file_path, skiprows=1, header=None,  delimiter=" ", usecols=range(1, None))
  # # Combinando as DataFrames
  # df = pd.concat([df_fixa, df_variavel], axis=1)
  # # Definindo nomes de colunas (opcional)
  # df.columns = ["dateTime", "sequential", "response_time", ...]
  # # Imprimindo o DataFrame
  # print(df.to_string())

  # print(df)
  df = df.dropna(subset=['sequential', 'response_time'])
  # print(df)
  df["dateTime"] = pd.to_datetime(df["dateTime"], format="mixed")
  df['duration'] = (df['dateTime'] - df['dateTime'].iloc[0]).dt.total_seconds()
  return df

def generate_boxplots(df, experiment, app, metric, level):
  """
  Gera boxplots para a métrica especificada comparando os diferentes protocolos.

  Args:
    df: DataFrame do Pandas com os dados do experimento.
    metric: Métrica a ser comparada ("memory" ou "cpu").
    app: "client" ou "server".
    level: experiment level

  Returns:
    Figura do Matplotlib com os boxplots.
  """
  if df.empty:
    return
  fig, ax = plt.subplots()
  fig.set_size_inches(14, 8)
  metric_column = "memory_usage(%)" if metric == "memory" else "cpu_usage(%)"
  # df[["dateTime", "duration", "protocol", "memory_usage(%)", "cpu_usage(%)"]].to_csv("df.csv")
  sns.boxplot(x="protocol", y=metric_column, data=df, ax=ax)
  ax.set_xlabel("Protocolo")
  ax.set_ylabel("% Memória Utilizada" if metric == "memory" else "% CPU Utilizado")
  ax.set_title(f"{experiment.capitalize()} - {app.capitalize()} - {metric.capitalize()} - {level}")
  # plt.show()
  return fig

def generate_lineplots_by_metric(df, experiment, app, metric, level):
  """
  Gera lineplots para a métrica especificada comparando os diferentes protocolos.

  Args:
    df: DataFrame do Pandas com os dados do experimento.
    metric: Métrica a ser comparada ("memory" ou "cpu").
    app: "client" ou "server".
    level: experiment level

  Returns:
    Figura do Matplotlib com os lineplots.
  """
  if df.empty:
    return
  fig, ax = plt.subplots()
  fig.set_size_inches(18, 8)
  metric_column = "memory_usage(%)" if metric == "memory" else "cpu_usage(%)"
  # df[["dateTime", "duration", "protocol", "memory_usage(%)", "cpu_usage(%)"]].to_csv("df.csv")
  sns.lineplot(x="duration", y=metric_column, data=df, hue="protocol")
  ax.set_xlabel("Duração (s)")
  ax.set_ylabel("% Memória Utilizada" if metric == "memory" else "% CPU Utilizado")
  ax.set_title(f"{experiment.capitalize()} - {app.capitalize()} - {metric.capitalize()} - {level}")
  plt.legend(bbox_to_anchor=(1.05, 1), loc='upper left')
  # plt.show()
  return fig

def generate_lineplots_by_response_time(df, experiment, level):
  """
  Gera lineplots para a métrica especificada comparando os diferentes protocolos.

  Args:
    df: DataFrame do Pandas com os dados do experimento.
    metric: Métrica a ser comparada ("memory" ou "cpu").
    app: "client" ou "server".
    level: experiment level

  Returns:
    Figura do Matplotlib com os lineplots.
  """
  if df.empty:
    return
  fig, ax = plt.subplots()
  fig.set_size_inches(14, 8)
  # sns.lineplot(x="duration", y="response_time", data=df, hue="protocol", sort=False)
  sns.boxplot(x="protocol", y="response_time", data=df, ax=ax)
  ax.set_xlabel("Duração (s)")
  ax.set_ylabel("Tempo de Resposta (ms)")
  # ax.set_ylim(bottom=0, top=max(df["response_time"]))
  ax.set_title(f"{experiment.capitalize()} - {level}")
  # plt.show()
  return fig


def save_plots(fig, output_directory, experiment, app, metric, level, kind):
  """
  Salva a figura do Matplotlib como um arquivo PNG.

  Args:
    fig: Figura do Matplotlib com os boxplots.
    output_directory: Caminho para a pasta onde os boxplots serão salvos.
    metric: Métrica a ser comparada ("memory" ou "cpu").
    app: "client" ou "server".
    level: experiment level
  """
  if fig is None:
    return
  if not os.path.exists(output_directory):
    os.makedirs(output_directory)
  file_name = f"{experiment}_{level}_{app}_{metric}_{kind}.png"
  fig.savefig(os.path.join(output_directory, file_name))

def main():
  """
  Função principal que gera os boxplots para os experimentos.
  """
  input_directory = "../results20240403-AllExecutedOk"
  output_directory = "./charts"

  experiments = ["Fibonacci", "SendFile"]
  fibonacci_levels = ["2", "11", "38"]
  sendfile_levels = ["sm", "md", "lg"]
  protocols = ["UDP", "TCP", "TLS", "RPC", "QUIC", "HTTP", "HTTPS", "HTTP2", "TCPTLS", "RPCHTTP", "TCPHTTP", "TLSHTTP2"]
  metrics = ["memory", "cpu"]
  apps = ["client", "server"]
  for experiment in experiments:
    levels = fibonacci_levels if experiment == "Fibonacci" else sendfile_levels
    for level in levels:
      for app in apps:
        for metric in metrics:
          df_monitor = pd.DataFrame()
          df_results = pd.DataFrame()
          for protocol in protocols:
            print("experiment/level/app/metric/protocol:", experiment, "/", level, "/", app, "/", metric, "/", protocol)
            for experiment_directory in os.listdir(input_directory):
              # print(experiment_directory, experiment_directory.upper())
              if experiment in experiment_directory and "-"+protocol+"-" in experiment_directory.upper() and "-"+level in experiment_directory:
                ############# Read Monitor Data #############
                file_path = get_last_log_file(os.path.join(input_directory, experiment_directory), app, "monitor")
                if file_path is None:
                  continue
                if not validate_file(file_path, "monitor"):
                  continue
                df_experiment = read_monitor_data(file_path)
                df_experiment["protocol"] = protocol
                #   df_monitor = df_monitor.append(df_experiment)
                df_concat = pd.concat([df_monitor, df_experiment], ignore_index=True)
                df_monitor = df_concat

                if df_monitor.empty:
                  continue

                ############# Read Results Data #############
                # avoid executing this block twice for the same experiment
                if metric == "memory" and app == "client":
                  file_path = get_last_log_file(os.path.join(input_directory, experiment_directory), app, "results")
                  if file_path is None:
                    continue
                  if not validate_file(file_path, "results"):
                    continue
                  df_experiment = read_results_data(file_path)
                  df_experiment["protocol"] = protocol
                  df_concat = pd.concat([df_results, df_experiment], ignore_index=True)
                  df_results = df_concat

                  if df_results.empty:
                    continue

              # df["duration"] = calculate_duration(df)
              # df = df[df["duration"] > 0]

          fig = generate_boxplots(df_monitor, experiment, app, metric, level)
          save_plots(fig, output_directory, experiment, app, metric, level, kind="boxplot")
          plt.close()
          fig = generate_lineplots_by_metric(df_monitor, experiment, app, metric, level)
          save_plots(fig, output_directory, experiment, app, metric, level, kind="lineplot")
          plt.close()
          if metric == "memory" and app == "client":
            fig = generate_lineplots_by_response_time(df_results, experiment, level)
            save_plots(fig, output_directory, experiment, app, "responsetime", level, kind="boxplot")
            plt.close()


if __name__ == "__main__":
    main()