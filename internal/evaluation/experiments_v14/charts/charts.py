import os
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

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
  kind_filter = f".{kind}." + "csv" if kind == "monitor" else "txt"
  log_files = [f for f in files if f.endswith(kind_filter) and app in f]
  if not log_files:
    return None
  return os.path.join(directory, max(log_files))

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
  header = ["dateTime", "sequential", "response_time"]
  print(file_path)
  df = pd.read_csv(file_path, header=None, delimiter=" ", skiprows=1, names=header)
  df = df.dropna(subset=['sequential', 'response_time'])
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
  fig, ax = plt.subplots()
  metric_column = "memory_usage(%)" if metric == "memory" else "cpu_usage(%)"
  # df[["dateTime", "duration", "protocol", "memory_usage(%)", "cpu_usage(%)"]].to_csv("df.csv")
  sns.boxplot(x="protocol", y=metric_column, data=df, ax=ax) #df[df["container_status"] == client_or_server]
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
  fig, ax = plt.subplots()
  metric_column = "memory_usage(%)" if metric == "memory" else "cpu_usage(%)"
  # df[["dateTime", "duration", "protocol", "memory_usage(%)", "cpu_usage(%)"]].to_csv("df.csv")
  sns.lineplot(x="duration", y=metric_column, data=df, hue="protocol") #df[df["container_status"] == client_or_server]
  ax.set_xlabel("Duração (s)")
  ax.set_ylabel("% Memória Utilizada" if metric == "memory" else "% CPU Utilizado")
  ax.set_title(f"{experiment.capitalize()} - {app.capitalize()} - {metric.capitalize()} - {level}")
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
  fig, ax = plt.subplots()
  sns.lineplot(x="duration", y="response_time", data=df, hue="protocol")
  ax.set_xlabel("Duração (s)")
  ax.set_ylabel("Tempo de Resposta (ms)")
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
  if not os.path.exists(output_directory):
    os.makedirs(output_directory)
  file_name = f"{experiment}_{app}_{metric}_{level}_{kind}.png"
  fig.savefig(os.path.join(output_directory, file_name))

def main():
  """
  Função principal que gera os boxplots para os experimentos.
  """
  input_directory = "../results"
  output_directory = "./charts"

  experiments = ["Fibonacci", "SendFile"]
  fibonacci_levels = ["2", "11", "38"]
  sendfile_levels = ["sm", "md", "lg"]
  protocols = ["UDP", "TCP", "TLS", "RPC", "QUIC", "HTTP", "HTTPS", "HTTP2"]
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
              if experiment in experiment_directory and protocol in experiment_directory.upper() and level in experiment_directory:
                ############# Read Monitor Data #############
                file_path = get_last_log_file(os.path.join(input_directory, experiment_directory), app, "monitor")
                if file_path is None:
                  continue
                df_experiment = read_monitor_data(file_path)
                df_experiment["protocol"] = protocol
                #   df_monitor = df_monitor.append(df_experiment)
                df_concat = pd.concat([df_monitor, df_experiment], ignore_index=True)
                df_monitor = df_concat

                if df_monitor.empty:
                  continue

                ############# Read Results Data #############
                file_path = get_last_log_file(os.path.join(input_directory, experiment_directory), app, "results")
                if file_path is None:
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
          save_plots(fig, output_directory, experiment, app, metric, level, kind="lineplot_by_metric")
          plt.close()
          fig = generate_lineplots_by_response_time(df_results, experiment, level)
          save_plots(fig, output_directory, experiment, app, metric, level, kind="lineplot_by_response_time")
          plt.close()

if __name__ == "__main__":
    main()