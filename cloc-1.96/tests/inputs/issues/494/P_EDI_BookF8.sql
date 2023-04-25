USE 
GO
/****** Object:  StoredProcedure [dbo].[P_EDI_Book]    Script Date: 2020/05/22 11:41:27 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO


--=================================================

---------------------------------------------------
--  
---------------------------------------------------

--  ＜下り＞
--
   ,@ext_day int = -5						--2017-02-28 更新
--   ,@ext_day int = -1						--2016-08-17 更新
   ,@ext_time datetime = '00:00:00'
   
AS
--=================================================
-- 開始
--=================================================
SP_MAIN:
  SET NOCOUNT ON

--=============================================
-- なんとなく共通処理
--=============================================

  -- 必要なローカル変数定義
  Declare @tran_ct int           -- トランザクションカウント
 
  Declare @FETCH_STATUS int -- @@FETCH_STATUSワーク

--=================================================
-- パラメタチェック
--=================================================

  Declare @a_tax_rate_type char(1)  -- 税端数処理区分
  -- ▼▼▼ 2013/11/21 add
  Declare @a_tonets_v_up_flag char(1)	-- トーハンTONET=V上り実装フラグ
  Declare @a_tonets_v_dw_flag char(1)	-- トーハンTONET=V下り実装フラグ
  -- ▲▲▲ 2013/11/21 add

--=================================================
-- 共通処理
--=================================================

  -- 下り用処理カウンタ
  Declare @count_all int  -- 処理対象
 

  -- サーバ名・データベース名セット
  Declare @server_name varchar(100)

  End

  -- ▼▼▼ 2013/11/21 add
  -- 上り在庫変動処理（TONETS-V:T1）
  If @edi_type = '56U' Begin

  End
  -- ▲▲▲ 2013/11/21 add
  
  -- 【上り未登録商品マスタ(エコール)】上り未登録商品マスタのクライアント処理	-- 2014/01/20 add
  If @edi_type = '60U' Begin

  -- ▼▼▼ 2014/01/28 add
  -- 【JEUGIA基幹連動】上り売上処理
  If @edi_type = '62U' Begin

  -- ▲▲▲ 2014/01/28 add

  --

--=================================================
-- SP終了
--=================================================

